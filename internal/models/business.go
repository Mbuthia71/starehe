package models

import "time"

type BusinessListing struct {
	ID              string     `json:"id" db:"id"`
	UserID          string     `json:"user_id" db:"user_id"`
	BusinessName    string     `json:"business_name" db:"business_name"`
	PhoneNumber     string     `json:"phone_number" db:"phone_number"`
	Location        string     `json:"location" db:"location"`
	Website         *string    `json:"website,omitempty" db:"website"`
	Description     string     `json:"description" db:"description"`
	InstagramHandle *string    `json:"instagram_handle,omitempty" db:"instagram_handle"`
	FacebookHandle  *string    `json:"facebook_handle,omitempty" db:"facebook_handle"`
	LogoURL         *string    `json:"logo_url,omitempty" db:"logo_url"`
	BannerURL       *string    `json:"banner_url,omitempty" db:"banner_url"`
	IsVerified      bool       `json:"is_verified" db:"is_verified"`
	IsFeatured      bool       `json:"is_featured" db:"is_featured"`
	FeaturedUntil   *time.Time `json:"featured_until,omitempty" db:"featured_until"`
	Status          string     `json:"status" db:"status"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

type Job struct {
	ID                  string     `json:"id" db:"id"`
	UserID              string     `json:"user_id" db:"user_id"`
	BusinessID          *string    `json:"business_id,omitempty" db:"business_id"`
	Title               string     `json:"title" db:"title"`
	Description         string     `json:"description" db:"description"`
	Requirements        *string    `json:"requirements,omitempty" db:"requirements"`
	Responsibilities    *string    `json:"responsibilities,omitempty" db:"responsibilities"`
	Location            *string    `json:"location,omitempty" db:"location"`
	JobType             string     `json:"job_type" db:"job_type"`
	SalaryRange         *string    `json:"salary_range,omitempty" db:"salary_range"`
	ApplicationDeadline *time.Time `json:"application_deadline,omitempty" db:"application_deadline"`
	Status              string     `json:"status" db:"status"`
	ViewsCount          int        `json:"views_count" db:"views_count"`
	ApplicationsCount   int        `json:"applications_count" db:"applications_count"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`
}

type JobApplication struct {
	ID           string    `json:"id" db:"id"`
	JobID        string    `json:"job_id" db:"job_id"`
	UserID       string    `json:"user_id" db:"user_id"`
	CoverLetter  *string   `json:"cover_letter,omitempty" db:"cover_letter"`
	ResumeURL    *string   `json:"resume_url,omitempty" db:"resume_url"`
	Status       string    `json:"status" db:"status"`
	AppliedAt    time.Time `json:"applied_at" db:"applied_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Tender struct {
	ID              string     `json:"id" db:"id"`
	UserID          string     `json:"user_id" db:"user_id"`
	OrganizationName string    `json:"organization_name" db:"organization_name"`
	Title           string     `json:"title" db:"title"`
	Description     string     `json:"description" db:"description"`
	Requirements    *string    `json:"requirements,omitempty" db:"requirements"`
	BudgetRange     *string    `json:"budget_range,omitempty" db:"budget_range"`
	SubmissionDeadline time.Time `json:"submission_deadline" db:"submission_deadline"`
	TenderNumber    *string    `json:"tender_number,omitempty" db:"tender_number"`
	Category        *string    `json:"category,omitempty" db:"category"`
	Status          string     `json:"status" db:"status"`
	ViewsCount      int        `json:"views_count" db:"views_count"`
	BidsCount       int        `json:"bids_count" db:"bids_count"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

type TenderBid struct {
	ID                  string    `json:"id" db:"id"`
	TenderID            string    `json:"tender_id" db:"tender_id"`
	UserID              string    `json:"user_id" db:"user_id"`
	CompanyName         string    `json:"company_name" db:"company_name"`
	ProposalDescription string    `json:"proposal_description" db:"proposal_description"`
	BidAmount           *float64  `json:"bid_amount,omitempty" db:"bid_amount"`
	ProposalURL         *string   `json:"proposal_url,omitempty" db:"proposal_url"`
	Status              string    `json:"status" db:"status"`
	SubmittedAt         time.Time `json:"submitted_at" db:"submitted_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

type ClassGroup struct {
	ID            string     `json:"id" db:"id"`
	SchoolType    string     `json:"school_type" db:"school_type"`
	YearOfCompletion int      `json:"year_of_completion" db:"year_of_completion"`
	ClassName     string     `json:"class_name" db:"class_name"`
	Description   *string    `json:"description,omitempty" db:"description"`
	ClassRepID    *string    `json:"class_rep_id,omitempty" db:"class_rep_id"`
	MemberCount   int        `json:"member_count" db:"member_count"`
	IsActive      bool       `json:"is_active" db:"is_active"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

type ClassGroupMember struct {
	ID           string    `json:"id" db:"id"`
	ClassGroupID string    `json:"class_group_id" db:"class_group_id"`
	UserID       string    `json:"user_id" db:"user_id"`
	Role         string    `json:"role" db:"role"`
	JoinedAt     time.Time `json:"joined_at" db:"joined_at"`
}

type MerchantOffer struct {
	ID                string     `json:"id" db:"id"`
	UserID            string     `json:"user_id" db:"user_id"`
	BusinessID        *string    `json:"business_id,omitempty" db:"business_id"`
	Title             string     `json:"title" db:"title"`
	Description       string     `json:"description" db:"description"`
	DiscountPercentage *int      `json:"discount_percentage,omitempty" db:"discount_percentage"`
	OriginalPrice     *float64   `json:"original_price,omitempty" db:"original_price"`
	OfferPrice        *float64   `json:"offer_price,omitempty" db:"offer_price"`
	ValidFrom         time.Time  `json:"valid_from" db:"valid_from"`
	ValidUntil        time.Time  `json:"valid_until" db:"valid_until"`
	TermsConditions   *string    `json:"terms_conditions,omitempty" db:"terms_conditions"`
	ImageURL          *string    `json:"image_url,omitempty" db:"image_url"`
	IsExclusive       bool       `json:"is_exclusive" db:"is_exclusive"`
	Status            string     `json:"status" db:"status"`
	ViewsCount        int        `json:"views_count" db:"views_count"`
	RedemptionsCount  int        `json:"redemptions_count" db:"redemptions_count"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}

type Sponsorship struct {
	ID            string     `json:"id" db:"id"`
	UserID        string     `json:"user_id" db:"user_id"`
	Title         string     `json:"title" db:"title"`
	Description   string     `json:"description" db:"description"`
	SponsorshipType string    `json:"sponsorship_type" db:"sponsorship_type"`
	TargetAmount  float64    `json:"target_amount" db:"target_amount"`
	CurrentAmount float64    `json:"current_amount" db:"current_amount"`
	StartDate     time.Time  `json:"start_date" db:"start_date"`
	EndDate       *time.Time `json:"end_date,omitempty" db:"end_date"`
	Status        string     `json:"status" db:"status"`
	Beneficiary   *string    `json:"beneficiary,omitempty" db:"beneficiary"`
	ImageURL      *string    `json:"image_url,omitempty" db:"image_url"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

type SponsorshipContribution struct {
	ID             string    `json:"id" db:"id"`
	SponsorshipID  string    `json:"sponsorship_id" db:"sponsorship_id"`
	UserID         string    `json:"user_id" db:"user_id"`
	Amount         float64   `json:"amount" db:"amount"`
	Message        *string   `json:"message,omitempty" db:"message"`
	IsAnonymous    bool      `json:"is_anonymous" db:"is_anonymous"`
	ContributedAt  time.Time `json:"contributed_at" db:"contributed_at"`
}

type EscrowTransaction struct {
	ID                string     `json:"id" db:"id"`
	BusinessID        string     `json:"business_id" db:"business_id"`
	BuyerID           string     `json:"buyer_id" db:"buyer_id"`
	SellerID          string     `json:"seller_id" db:"seller_id"`
	Amount            float64    `json:"amount" db:"amount"`
	Description       string     `json:"description" db:"description"`
	Status            string     `json:"status" db:"status"`
	ReleaseConditions *string    `json:"release_conditions,omitempty" db:"release_conditions"`
	FundedAt          *time.Time `json:"funded_at,omitempty" db:"funded_at"`
	CompletedAt       *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	CancelledAt       *time.Time `json:"cancelled_at,omitempty" db:"cancelled_at"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at" db:"updated_at"`
}
