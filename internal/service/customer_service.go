package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/nanwannaporn/insurance-system/internal/domain"
	"gorm.io/gorm"
)

// Databas layer
type CustomerService struct {
	DB *gorm.DB
}

func NewCustomerService(db *gorm.DB) *CustomerService {
	return &CustomerService{DB: db}
}

func (s *CustomerService) CreateCustomer(c *domain.Customer) (string, error) {
	if c.Age < 18 {
		return "", errors.New("customer must be at least 18 years old to proceed")
	}

	newID := fmt.Sprintf("C-%d", time.Now().UnixNano())
	c.CustomerID = newID

	result := s.DB.Create(c)

	if result.Error != nil {
		return "", fmt.Errorf("database error: failed to insert customer: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return "", errors.New("customer record was not created; zero rows affected")
	}

	return newID, nil
}

func (s *CustomerService) UpdateBeneficiaries(customerID string, beneficiaries []domain.Beneficiaries) error {
	var customer domain.Customer
	result := s.DB.Where("customer_id = ? and deleted_at is null", customerID).First(&customer)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("customer not found")
		}
		return fmt.Errorf("database query error during customer check: %w", result.Error)
	}

	delete_result := s.DB.Where("customer_id = ?", customerID).Delete(&domain.Beneficiaries{})
	if delete_result.Error != nil {
		return fmt.Errorf("failed to delete old beneficiaries: %w", delete_result.Error)
	}

	for i := range beneficiaries {
		beneficiaries[i].CustomerID = customerID
	}

	create_result := s.DB.Create(&beneficiaries)
	if create_result.Error != nil {
		return fmt.Errorf("failed to create new beneficiaries: %w", create_result.Error)
	}

	return nil

}

func (s *CustomerService) UpdateHealthDeclaration(customerID string, health *domain.HealthDeclaration) error {
	var customer domain.Customer
	if err := s.DB.Where("customer_id = ? and deleted_at is null", customerID).First(&customer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("customer not found")
		}
		return fmt.Errorf("database query error: %w", err)
	}

	health.CustomerID = customerID

	var existing domain.HealthDeclaration
	result := s.DB.Where("customer_id = ?", customerID).First(&existing)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check existing health record: %w", result.Error)
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		if err := s.DB.Create(health).Error; err != nil {
			return fmt.Errorf("failed to create new health record: %w", err)
		}
	} else {
		health.ID = existing.ID
		if err := s.DB.Save(health).Error; err != nil {
			return fmt.Errorf("failed to update existing health record: %w", err)
		}
	}

	return nil
}

func (s *CustomerService) GetPlans(age int, sumAssured float64) ([]domain.Insurance, error) {
	var plans []domain.Insurance

	tx := s.DB
	if age > 0 && sumAssured > 0 {
		tx = tx.Where("? >= min_age and ? <= max_age and sum_assured_limit <= ?", age, age, sumAssured)
	} else if age > 0 {
		tx = tx.Where("? >= min_age and ? <= max_age", age, age)
	} else if sumAssured > 0 {
		tx = tx.Where("sum_assured_limit <= ?", sumAssured)
	}

	if err := tx.Find(&plans).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve insurance plans: %w", err)
	}

	return plans, nil
}

func (s *CustomerService) PurchaseInsurance(req *domain.InsurancePurchaseRequest) (*domain.CustomerInsurance, error) {
	var customer domain.Customer
	if err := s.DB.Where("customer_id = ? and deleted_at is null", req.CustomerID).First(&customer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("customer not found")
		}
		return nil, fmt.Errorf("db error checking customer: %w", err)
	}

	var plan domain.Insurance
	if err := s.DB.Where("id = ? and deleted_at is null", req.InsuranceID).First(&plan).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("plan not found")
		}
		return nil, fmt.Errorf("db error checking plan: %w", err)
	}

	var ownInsurance domain.CustomerInsurance
	result := s.DB.Where("customer_id = ? and insurance_id = ? and (status = ? or status = ?)", req.CustomerID, req.InsuranceID, "Active", "Pending").First(&ownInsurance)
	if result.Error == nil {
		return nil, errors.New("customer already owns an active policy with this Insurance ID")
	}

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// Handle other database errors besides "not found"
		return nil, fmt.Errorf("database error checking existing policy: %w", result.Error)
	}

	customerSuffix := req.CustomerID[len(req.CustomerID)-4:]
	policyNumber := fmt.Sprintf("P%d-%s-%d", req.InsuranceID, customerSuffix, time.Now().Unix())
	fmt.Println("policynum", policyNumber)

	insurancePolicy := domain.CustomerInsurance{
		PolicyNumber:     policyNumber,
		CustomerID:       req.CustomerID,
		InsuranceID:      req.InsuranceID,
		SumAssured:       plan.SumAssuredLimit,
		PremiumAmount:    plan.InsurancePremium,
		PaymentFrequency: req.PaymentFrequency,
		PaymentMethod:    req.PaymentMethod,
		EffectDate:       time.Now().AddDate(0, 0, 15),
		Status:           "Pending",
	}

	tx := s.DB.Begin()

	if err := s.DB.Create(&insurancePolicy).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create policy record: %w", err)
	}

	//update Beneficiaries
	var beneficiary domain.Beneficiaries
	updateResult := s.DB.Model(&beneficiary).Where("customer_id = ? and insurance_id = ?", req.CustomerID, req.InsuranceID).Update("PolicyNumber", policyNumber)
	if updateResult.Error != nil {
		return nil, fmt.Errorf("failed to link policy number to beneficiaries: %w", updateResult.Error)
	}

	return &insurancePolicy, nil

}
