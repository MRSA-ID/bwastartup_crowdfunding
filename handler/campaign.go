package handler

import (
	"bwastartup/campaign"
	"bwastartup/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// tangkap parameter di handler
//  handler ke service
// service yang menentukan repository mana yang di-call
// repository : Get All, GetByUSerID
// db

type campaignHandler struct{
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler{
	return &campaignHandler{service}
}

// api/v1/campaigns
func(h *campaignHandler) GetCampaigns(c *gin.Context){
	userID, _:= strconv.Atoi(c.Query("user_id"))
	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil{
		response := helper.APIResponse("Error to get campaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
}

func (h *campaignHandler) GetCampaign(c *gin.Context){
	// api/v1/campaign/1
	// handler : mapping id yang di url ke struct input => service, formatter
	// service : input nya struct input => menangkap id di URL, manggil repo
	// repository : get campaign by ID

	var input campaign.GetCampaignDetailInput

	err := c.ShouldBindUri(&input)

	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	
	campaignDetail, err := h.service.GetCampaignByID(input)

	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Campaign detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK,response)
}