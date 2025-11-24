package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nanwannaporn/insurance-system/internal/domain"
)

type InsuranceService interface {
	CreateCustomer(c *domain.Customer) (string, error)
	UpdateBeneficiaries(customerID string, beneficiaries []domain.Beneficiaries) error
	UpdateHealthDeclaration(customerID string, health *domain.HealthDeclaration) error
	GetPlans(age int, sumAssured float64) ([]domain.Insurance, error)
	PurchaseInsurance(req *domain.InsurancePurchaseRequest) (*domain.CustomerInsurance, error)
}

type CustomerHandler struct {
	Service InsuranceService
}

func NewCustomerHandler(s InsuranceService) *CustomerHandler {
	return &CustomerHandler{
		Service: s,
	}
}

func (h *CustomerHandler) CreateCustomerHandler(c *gin.Context) {
	customer := new(domain.Customer)
	if err := c.BindJSON(customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Invalid input",
			"detail": err.Error(),
		})
		return
	}

	newID, err := h.Service.CreateCustomer(customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create customer",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Customer created successfully.",
		"customer_id": newID,
	})
}

func (h *CustomerHandler) UpdateBeneficiaryHandler(c *gin.Context) {
	customerID := c.Param("id")
	if customerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer ID is required"})
		return
	}

	var beneficiaries []domain.Beneficiaries
	if err := c.ShouldBindJSON(&beneficiaries); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid beneficiary info format or missing fields"})
		return
	}

	err := h.Service.UpdateBeneficiaries(customerID, beneficiaries)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update beneficiary information", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Beneficiary information updated successfully for Customer ID: %s", customerID),
	})

}

func (h *CustomerHandler) UpdateHealthDeclarationHandler(c *gin.Context) {
	customerID := c.Param("id")
	if customerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer ID is required in URL path"})
		return
	}

	var health domain.HealthDeclaration
	if err := c.ShouldBindJSON(&health); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid health info input", "details": err.Error()})
		return
	}

	err := h.Service.UpdateHealthDeclaration(customerID, &health)
	if err != nil {
		if err.Error() == "customer not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found, cannot update health data"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save health declaration", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Health declaration saved/updated successfully for Customer ID: %s", customerID),
	})
}

func (h *CustomerHandler) GetPlansHandler(c *gin.Context) {
	var age int
	var sumAssured float64

	ageStr := c.Query("age")
	if ageStr != "" {
		if a, err := strconv.Atoi(ageStr); err == nil {
			age = a
		}
	}

	sumAssuredStr := c.Query("sum_assured")
	if sumAssuredStr != "" {
		if s, err := strconv.ParseFloat(sumAssuredStr, 64); err == nil {
			sumAssured = s
		}
	}

	plans, err := h.Service.GetPlans(age, sumAssured)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch insurance plans",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, plans)
}

func (h *CustomerHandler) CreatePolicyHandler(c *gin.Context) {
	var req domain.InsurancePurchaseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format or missing fields", "details": err.Error()})
		return
	}

	customerInsurance, err := h.Service.PurchaseInsurance(&req)

	if err != nil {
		if err.Error() == "customer not found" || err.Error() == "plan not found" || strings.Contains(err.Error(), "exceeds plan limit") {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Policy creation failed", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":           "Policy issued successfully.",
		"customer_id":       customerInsurance.CustomerID,
		"policy_number":     customerInsurance.PolicyNumber,
		"sum_assured":       customerInsurance.SumAssured,
		"premium_amount":    customerInsurance.PremiumAmount,
		"payment_frequency": customerInsurance.PaymentFrequency,
		"payment_Method ":   customerInsurance.PaymentMethod,
		"status":            customerInsurance.Status,
	})

}
