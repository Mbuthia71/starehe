package business

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"starehian-society-platform/internal/services"
)

type BusinessHandler struct {
	businessService *services.BusinessService
}

func NewBusinessHandler(businessService *services.BusinessService) *BusinessHandler {
	return &BusinessHandler{
		businessService: businessService,
	}
}

// Business Listings
func (h *BusinessHandler) CreateBusinessListing(c *fiber.Ctx) error {
	var req struct {
		BusinessName    string  `json:"business_name"`
		PhoneNumber     string  `json:"phone_number"`
		Location        string  `json:"location"`
		Website         *string `json:"website"`
		Description     string  `json:"description"`
		InstagramHandle *string `json:"instagram_handle"`
		FacebookHandle  *string `json:"facebook_handle"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	userID := c.Locals("user_id").(string)
	business, err := h.businessService.CreateBusinessListing(c.Context(), userID, req.BusinessName, req.PhoneNumber, req.Location, req.Description, req.Website, req.InstagramHandle, req.FacebookHandle)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(business)
}

func (h *BusinessHandler) GetBusinessListing(c *fiber.Ctx) error {
	id := c.Params("id")
	business, err := h.businessService.GetBusinessListing(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if business == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Business listing not found"})
	}

	return c.JSON(business)
}

func (h *BusinessHandler) GetUserBusinessListings(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	listings, err := h.businessService.GetUserBusinessListings(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(listings)
}

func (h *BusinessHandler) ListBusinessListings(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	listings, err := h.businessService.ListBusinessListings(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(listings)
}

func (h *BusinessHandler) UpdateBusinessListing(c *fiber.Ctx) error {
	id := c.Params("id")
	var updates map[string]interface{}

	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	err := h.businessService.UpdateBusinessListing(c.Context(), id, updates)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Business listing updated successfully"})
}

func (h *BusinessHandler) DeleteBusinessListing(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.businessService.DeleteBusinessListing(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Business listing deleted successfully"})
}

// Jobs
func (h *BusinessHandler) CreateJob(c *fiber.Ctx) error {
	var req struct {
		Title               string  `json:"title"`
		Description         string  `json:"description"`
		BusinessID          *string `json:"business_id"`
		Requirements        *string `json:"requirements"`
		Responsibilities    *string `json:"responsibilities"`
		Location            *string `json:"location"`
		JobType             *string `json:"job_type"`
		SalaryRange         *string `json:"salary_range"`
		ApplicationDeadline *string `json:"application_deadline"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	userID := c.Locals("user_id").(string)
	job, err := h.businessService.CreateJob(c.Context(), userID, req.Title, req.Description, req.BusinessID, req.Requirements, req.Responsibilities, req.Location, req.JobType, req.SalaryRange, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(job)
}

func (h *BusinessHandler) GetJob(c *fiber.Ctx) error {
	id := c.Params("id")
	job, err := h.businessService.GetJob(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if job == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Job not found"})
	}

	return c.JSON(job)
}

func (h *BusinessHandler) ListJobs(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	jobs, err := h.businessService.ListJobs(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(jobs)
}

func (h *BusinessHandler) ApplyForJob(c *fiber.Ctx) error {
	jobID := c.Params("id")
	var req struct {
		CoverLetter *string `json:"cover_letter"`
		ResumeURL   *string `json:"resume_url"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	userID := c.Locals("user_id").(string)
	application, err := h.businessService.ApplyForJob(c.Context(), jobID, userID, req.CoverLetter, req.ResumeURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(application)
}

func (h *BusinessHandler) GetJobApplications(c *fiber.Ctx) error {
	jobID := c.Params("id")
	applications, err := h.businessService.GetJobApplications(c.Context(), jobID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(applications)
}

// Tenders
func (h *BusinessHandler) CreateTender(c *fiber.Ctx) error {
	var req struct {
		OrganizationName   string  `json:"organization_name"`
		Title              string  `json:"title"`
		Description        string  `json:"description"`
		Requirements       *string `json:"requirements"`
		BudgetRange        *string `json:"budget_range"`
		TenderNumber       *string `json:"tender_number"`
		Category           *string `json:"category"`
		SubmissionDeadline string  `json:"submission_deadline"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	userID := c.Locals("user_id").(string)
	tender, err := h.businessService.CreateTender(c.Context(), userID, req.OrganizationName, req.Title, req.Description, req.Requirements, req.BudgetRange, req.TenderNumber, req.Category, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(tender)
}

func (h *BusinessHandler) GetTender(c *fiber.Ctx) error {
	id := c.Params("id")
	tender, err := h.businessService.GetTender(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if tender == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Tender not found"})
	}

	return c.JSON(tender)
}

func (h *BusinessHandler) ListTenders(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	tenders, err := h.businessService.ListTenders(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(tenders)
}

func (h *BusinessHandler) SubmitTenderBid(c *fiber.Ctx) error {
	tenderID := c.Params("id")
	var req struct {
		CompanyName         string   `json:"company_name"`
		ProposalDescription string   `json:"proposal_description"`
		BidAmount           *float64 `json:"bid_amount"`
		ProposalURL         *string  `json:"proposal_url"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	userID := c.Locals("user_id").(string)
	bid, err := h.businessService.SubmitTenderBid(c.Context(), tenderID, userID, req.CompanyName, req.ProposalDescription, req.BidAmount, req.ProposalURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(bid)
}

func (h *BusinessHandler) GetTenderBids(c *fiber.Ctx) error {
	tenderID := c.Params("id")
	bids, err := h.businessService.GetTenderBids(c.Context(), tenderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(bids)
}

// Class Groups
func (h *BusinessHandler) CreateClassGroup(c *fiber.Ctx) error {
	var req struct {
		SchoolType      string  `json:"school_type"`
		YearOfCompletion int    `json:"year_of_completion"`
		ClassName       string  `json:"class_name"`
		Description     *string `json:"description"`
		ClassRepID      *string `json:"class_rep_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	group, err := h.businessService.CreateClassGroup(c.Context(), req.SchoolType, req.YearOfCompletion, req.ClassName, req.Description, req.ClassRepID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(group)
}

func (h *BusinessHandler) GetClassGroup(c *fiber.Ctx) error {
	id := c.Params("id")
	group, err := h.businessService.GetClassGroup(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if group == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Class group not found"})
	}

	return c.JSON(group)
}

func (h *BusinessHandler) ListClassGroups(c *fiber.Ctx) error {
	schoolType := c.Query("school_type")
	year, _ := strconv.Atoi(c.Query("year", "0"))

	groups, err := h.businessService.ListClassGroups(c.Context(), schoolType, year)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(groups)
}

func (h *BusinessHandler) JoinClassGroup(c *fiber.Ctx) error {
	groupID := c.Params("id")
	var req struct {
		Role string `json:"role"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	userID := c.Locals("user_id").(string)
	err := h.businessService.JoinClassGroup(c.Context(), groupID, userID, req.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Joined class group successfully"})
}

func (h *BusinessHandler) GetClassGroupMembers(c *fiber.Ctx) error {
	groupID := c.Params("id")
	members, err := h.businessService.GetClassGroupMembers(c.Context(), groupID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(members)
}

func (h *BusinessHandler) LeaveClassGroup(c *fiber.Ctx) error {
	groupID := c.Params("id")
	userID := c.Locals("user_id").(string)
	err := h.businessService.LeaveClassGroup(c.Context(), groupID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Left class group successfully"})
}

// Merchant Offers
func (h *BusinessHandler) CreateMerchantOffer(c *fiber.Ctx) error {
	var req struct {
		Title             string   `json:"title"`
		Description       string   `json:"description"`
		BusinessID        *string  `json:"business_id"`
		DiscountPercentage *int    `json:"discount_percentage"`
		OriginalPrice     *float64 `json:"original_price"`
		OfferPrice        *float64 `json:"offer_price"`
		ValidFrom         string   `json:"valid_from"`
		ValidUntil        string   `json:"valid_until"`
		TermsConditions   *string  `json:"terms_conditions"`
		ImageURL          *string  `json:"image_url"`
		IsExclusive       bool     `json:"is_exclusive"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	userID := c.Locals("user_id").(string)
	offer, err := h.businessService.CreateMerchantOffer(c.Context(), userID, req.Title, req.Description, req.BusinessID, req.DiscountPercentage, req.OriginalPrice, req.OfferPrice, nil, nil, req.TermsConditions, req.ImageURL, req.IsExclusive)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(offer)
}

func (h *BusinessHandler) GetMerchantOffer(c *fiber.Ctx) error {
	id := c.Params("id")
	offer, err := h.businessService.GetMerchantOffer(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if offer == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Merchant offer not found"})
	}

	return c.JSON(offer)
}

func (h *BusinessHandler) ListMerchantOffers(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	offers, err := h.businessService.ListMerchantOffers(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(offers)
}

// Sponsorships
func (h *BusinessHandler) CreateSponsorship(c *fiber.Ctx) error {
	var req struct {
		Title           string  `json:"title"`
		Description     string  `json:"description"`
		SponsorshipType string  `json:"sponsorship_type"`
		TargetAmount    float64 `json:"target_amount"`
		StartDate       string  `json:"start_date"`
		EndDate         *string `json:"end_date"`
		Beneficiary     *string `json:"beneficiary"`
		ImageURL        *string `json:"image_url"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	userID := c.Locals("user_id").(string)
	sponsorship, err := h.businessService.CreateSponsorship(c.Context(), userID, req.Title, req.Description, req.SponsorshipType, req.TargetAmount, nil, nil, req.Beneficiary, req.ImageURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(sponsorship)
}

func (h *BusinessHandler) GetSponsorship(c *fiber.Ctx) error {
	id := c.Params("id")
	sponsorship, err := h.businessService.GetSponsorship(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if sponsorship == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Sponsorship not found"})
	}

	return c.JSON(sponsorship)
}

func (h *BusinessHandler) ListSponsorships(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	sponsorships, err := h.businessService.ListSponsorships(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(sponsorships)
}

func (h *BusinessHandler) ContributeToSponsorship(c *fiber.Ctx) error {
	sponsorshipID := c.Params("id")
	var req struct {
		Amount      float64 `json:"amount"`
		Message     *string `json:"message"`
		IsAnonymous bool    `json:"is_anonymous"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	userID := c.Locals("user_id").(string)
	contribution, err := h.businessService.ContributeToSponsorship(c.Context(), sponsorshipID, userID, req.Amount, req.Message, req.IsAnonymous)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(contribution)
}

func (h *BusinessHandler) GetSponsorshipContributions(c *fiber.Ctx) error {
	sponsorshipID := c.Params("id")
	contributions, err := h.businessService.GetSponsorshipContributions(c.Context(), sponsorshipID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(contributions)
}

// Escrow Transactions
func (h *BusinessHandler) CreateEscrowTransaction(c *fiber.Ctx) error {
	var req struct {
		BusinessID        string  `json:"business_id"`
		SellerID          string  `json:"seller_id"`
		Amount            float64 `json:"amount"`
		Description       string  `json:"description"`
		ReleaseConditions *string `json:"release_conditions"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	userID := c.Locals("user_id").(string)
	transaction, err := h.businessService.CreateEscrowTransaction(c.Context(), req.BusinessID, userID, req.SellerID, req.Amount, req.Description, req.ReleaseConditions)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(transaction)
}

func (h *BusinessHandler) GetEscrowTransaction(c *fiber.Ctx) error {
	id := c.Params("id")
	transaction, err := h.businessService.GetEscrowTransaction(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if transaction == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Escrow transaction not found"})
	}

	return c.JSON(transaction)
}

func (h *BusinessHandler) GetUserEscrowTransactions(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	role := c.Query("role", "buyer")
	transactions, err := h.businessService.GetUserEscrowTransactions(c.Context(), userID, role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(transactions)
}

func (h *BusinessHandler) FundEscrowTransaction(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.businessService.FundEscrowTransaction(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Escrow transaction funded successfully"})
}

func (h *BusinessHandler) ReleaseEscrowTransaction(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.businessService.ReleaseEscrowTransaction(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Escrow transaction released successfully"})
}

func (h *BusinessHandler) CancelEscrowTransaction(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.businessService.CancelEscrowTransaction(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Escrow transaction cancelled successfully"})
}
