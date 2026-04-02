package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"rygell-dashboard/internal/services"

	"github.com/gin-gonic/gin"
)

type SearchHandler struct {
	masterService   *services.MasterService
	contractService *services.ContractService
}

type SearchItem struct {
	ID    uint   `json:"id"`
	Type  string `json:"type"`
	Label string `json:"label"`
}

func NewSearchHandler(masterService *services.MasterService, contractService *services.ContractService) *SearchHandler {
	return &SearchHandler{
		masterService:   masterService,
		contractService: contractService,
	}
}

func (h *SearchHandler) Search(c *gin.Context) {
	q := strings.TrimSpace(c.Query("q"))
	if len(q) < 2 {
		c.JSON(http.StatusOK, []SearchItem{})
		return
	}

	vendors, err := h.masterService.GetAllVendors(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	mills, err := h.masterService.GetAllMills(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	zones, err := h.masterService.GetAllZones(q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	emptyFilters := map[string]interface{}{}
	dedicatedFix, err := h.contractService.GetAllDedicatedFix(emptyFilters, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	dedicatedVar, err := h.contractService.GetAllDedicatedVar(emptyFilters, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	oncall, err := h.contractService.GetAllOncall(emptyFilters, q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	results := make([]SearchItem, 0, len(vendors)+len(mills)+len(zones)+len(dedicatedFix)+len(dedicatedVar)+len(oncall))

	for _, item := range vendors {
		results = append(results, SearchItem{
			ID:    item.ID,
			Type:  "vendors",
			Label: item.Name,
		})
	}
	for _, item := range mills {
		results = append(results, SearchItem{
			ID:    item.ID,
			Type:  "mills",
			Label: item.Name,
		})
	}
	for _, item := range zones {
		results = append(results, SearchItem{
			ID:    item.ID,
			Type:  "zones",
			Label: fmt.Sprintf("%s (%s)", item.Name, item.Type),
		})
	}
	for _, item := range dedicatedFix {
		results = append(results, SearchItem{
			ID:    item.ID,
			Type:  "contracts/dedicated-fix",
			Label: buildContractLabel(item.ID, item.SPKNumber, item.Vendor.Name, item.Mill.Name),
		})
	}
	for _, item := range dedicatedVar {
		originName := ""
		destName := ""
		if item.OriginZone != nil {
			originName = item.OriginZone.Name
		}
		if item.DestZone != nil {
			destName = item.DestZone.Name
		}
		results = append(results, SearchItem{
			ID:    item.ID,
			Type:  "contracts/dedicated-var",
			Label: buildContractLabelWithRoute(item.ID, item.SPKNumber, item.Vendor.Name, item.Mill.Name, originName, destName),
		})
	}
	for _, item := range oncall {
		originName := ""
		destName := ""
		if item.OriginZone != nil {
			originName = item.OriginZone.Name
		}
		if item.DestZone != nil {
			destName = item.DestZone.Name
		}
		results = append(results, SearchItem{
			ID:    item.ID,
			Type:  "contracts/oncall",
			Label: buildContractLabelWithRoute(item.ID, item.SPKNumber, item.Vendor.Name, item.Mill.Name, originName, destName),
		})
	}

	c.JSON(http.StatusOK, results)
}

func buildContractLabel(id uint, spkNumber string, vendorName string, millName string) string {
	base := strings.TrimSpace(spkNumber)
	if base == "" {
		base = fmt.Sprintf("Contract #%d", id)
	}

	detailParts := make([]string, 0, 2)
	if strings.TrimSpace(vendorName) != "" {
		detailParts = append(detailParts, strings.TrimSpace(vendorName))
	}
	if strings.TrimSpace(millName) != "" {
		detailParts = append(detailParts, strings.TrimSpace(millName))
	}
	if len(detailParts) == 0 {
		return base
	}

	return fmt.Sprintf("%s - %s", base, strings.Join(detailParts, " -> "))
}

func buildContractLabelWithRoute(id uint, spkNumber string, vendorName string, millName string, originZone string, destZone string) string {
	label := buildContractLabel(id, spkNumber, vendorName, millName)

	origin := strings.TrimSpace(originZone)
	dest := strings.TrimSpace(destZone)

	if origin != "" && dest != "" {
		label += fmt.Sprintf(" [%s → %s]", origin, dest)
	} else if origin != "" {
		label += fmt.Sprintf(" [%s]", origin)
	} else if dest != "" {
		label += fmt.Sprintf(" [→ %s]", dest)
	}

	return label
}
