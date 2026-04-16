package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"rygell-dashboard/internal/services"

	"github.com/gin-gonic/gin"
)

// searchLimit is the maximum number of results per category in global search.
const searchLimit = 10

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

	results := make([]SearchItem, 0, 50)

	// --- Master data searches (limited per category) ---

	if vendors, err := h.masterService.SearchVendorsLimited(q, searchLimit); err != nil {
		log.Printf("[SEARCH] vendors query failed for q=%q: %v", q, err)
	} else {
		for _, item := range vendors {
			results = append(results, SearchItem{
				ID:    item.ID,
				Type:  "vendors",
				Label: item.Name,
			})
		}
	}

	if mills, err := h.masterService.SearchMillsLimited(q, searchLimit); err != nil {
		log.Printf("[SEARCH] mills query failed for q=%q: %v", q, err)
	} else {
		for _, item := range mills {
			results = append(results, SearchItem{
				ID:    item.ID,
				Type:  "mills",
				Label: item.Name,
			})
		}
	}

	if zones, err := h.masterService.SearchZonesLimited(q, searchLimit); err != nil {
		log.Printf("[SEARCH] zones query failed for q=%q: %v", q, err)
	} else {
		for _, item := range zones {
			results = append(results, SearchItem{
				ID:    item.ID,
				Type:  "zones",
				Label: fmt.Sprintf("%s (%s)", item.Name, item.Type),
			})
		}
	}

	// --- Contract searches (using paginated queries with limit) ---

	emptyFilters := map[string]interface{}{}

	if dedicatedFix, _, err := h.contractService.GetDedicatedFixPage(emptyFilters, q, searchLimit, 0); err != nil {
		log.Printf("[SEARCH] dedicated-fix query failed for q=%q: %v", q, err)
	} else {
		for _, item := range dedicatedFix {
			results = append(results, SearchItem{
				ID:    item.ID,
				Type:  "contracts/dedicated-fix",
				Label: buildContractLabel(item.ID, item.SPKNumber, item.Vendor.Name, item.Mill.Name),
			})
		}
	}

	if dedicatedVar, _, err := h.contractService.GetDedicatedVarPage(emptyFilters, q, searchLimit, 0); err != nil {
		log.Printf("[SEARCH] dedicated-var query failed for q=%q: %v", q, err)
	} else {
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
	}

	if oncall, _, err := h.contractService.GetOncallPage(emptyFilters, q, searchLimit, 0); err != nil {
		log.Printf("[SEARCH] oncall query failed for q=%q: %v", q, err)
	} else {
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
	}

	log.Printf("[SEARCH] q=%q => %d results", q, len(results))
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
		label += fmt.Sprintf(" [%s -> %s]", origin, dest)
	} else if origin != "" {
		label += fmt.Sprintf(" [%s]", origin)
	} else if dest != "" {
		label += fmt.Sprintf(" [-> %s]", dest)
	}

	return label
}
