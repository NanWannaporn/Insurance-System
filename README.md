## 🛡️Insurance Management System (Go/GORM)

This project is a RESTful API backend designed to manage core insurance operations,
including customer records, beneficiary details, health declarations, and policy issuance. 
It leverages Go's performance and GORM for robust database interaction.

🚀 Key Features
- Customer Lifecycle Management: Create and update core customer details.
- Sub-Resource Updates: Dedicated endpoints for updating secondary data (Beneficiaries and Health Declarations) linked to a CustomerID.
- Dynamic Insurance Filtering: Retrieve available insurance plans based on dynamic criteria (Age and Sum Assured).
- Policy Issuance (Purchase): A dedicated transaction endpoint that handles business logic:
  - Verifies customer and plan existence.
  - Prevents duplicate active policies for the same plan.
  - Generates a unique PolicyNumber concatenated with the last 4 digits of the CustomerID.
  - Atomically links the new PolicyNumber back to the customer's Beneficiary records.
- Soft Deletes: Utilizes GORM's DeletedAt for non-destructive record deletion.


🛠️ Technology Stack
- Backend Language: Go
- Web Framework: Gin
- Database ORM: GORM
- Database: PostgreSQL

  
🏃 Getting Started 

1. Clone the Repository: git clone https://github.com/NanWannaporn/Insurance-System.git
2. Install Dependencies: go mod tidy
3. Run: go run ./cmd/main.go
