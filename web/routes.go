package web

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	icp2 "sysafari.com/customs/tguard/icp"
	"sysafari.com/customs/tguard/utils"
	"time"
)

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// AppendToICP
// @Summary      Add the specified Customs IDs to the specified ICP file
// @Description  If the specified ICP file does not exist will create a new ICP file with the specified ICP file name
// @Tags         icp
// @Accept       json
// @Produce      json
// @Param 		 message body CustomsAppendToICP true "Customs append into the ICP file"
// @Success      200
// @Failure      400
// @Router       /icp/append [post]
func AppendToICP(c echo.Context) (err error) {
	var errs []string
	aicp := new(CustomsAppendToICP)
	if err = c.Bind(aicp); err != nil {
		errs = append(errs, err.Error())
	}
	if err = c.Validate(aicp); err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return c.JSON(http.StatusBadRequest, &IcpResponse{
			Status: FAIL,
			Errors: errs,
		})
	}

	icp := &icp2.FileOfICP{
		FileName:   aicp.FileName,
		CustomsIDs: aicp.CustomsIds,
	}

	filename := icp.GenerateICP()
	errs = icp.Errors

	if errs != nil && len(errs) > 0 {
		return c.JSON(http.StatusInternalServerError, &IcpResponse{
			Status: FAIL,
			Errors: errs,
		})
	}

	return c.JSON(http.StatusOK, &IcpResponse{
		Status:   SUCCESS,
		FileName: filename,
	})
}

// MakeICPForTaxAgency
// @Summary      Generate a month's ICP file for tax agency
// @Description  If there is no customs declaration in the specified month of the tax agency, it will not be generated
// @Tags         icp
// @Accept       json
// @Produce      json
// @Param 		 dutyParty path string true "The duty party of tax agency"
// @Param 		 month query string false "which month, default is this month,example:2006-01"
// @Success      200
// @Failure      400
// @Router       /icp/taxAgency/{dutyParty} [get]
func MakeICPForTaxAgency(c echo.Context) (err error) {
	var errs []string
	dutyParty := c.Param("dutyParty")
	if dutyParty == "" {
		return c.JSON(http.StatusBadRequest, &IcpResponse{
			Status: FAIL,
			Errors: []string{fmt.Sprintf("The duty party is required.")},
		})
	}
	month := c.QueryParam("month")
	_, err = time.Parse("2006-01", month)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &IcpResponse{
			Status: FAIL,
			Errors: []string{fmt.Sprintf("The duty party is required.")},
		})
	}

	if month == "" {
		month = time.Now().Format("2006-01")
		log.Printf("Month is empty, use this month:%s instead.\n", month)
	}

	icp := &icp2.FileOfICP{
		DutyParty: dutyParty,
	}

	startDate := fmt.Sprintf("%s-01", month)
	endDate := fmt.Sprintf("%s-31", month)
	icp.QueryCustomsIDs(startDate, endDate)

	filename := icp.GenerateICP()

	errs = icp.Errors
	if errs != nil && len(errs) > 0 {
		return c.JSON(http.StatusInternalServerError, &IcpResponse{
			Status: FAIL,
			Errors: errs,
		})
	}

	return c.JSON(http.StatusOK, &IcpResponse{
		Status:   SUCCESS,
		FileName: filename,
	})
}

// DownloadFile
// Download ICP
// @Summary      Download ICP file
// @Description  File name format (BE0796544895_202209_01154020.xlsx), the file path will be found by the date in the file name
// @Tags         download
// @Accept       json
// @Produce      json
// @Param        filename   path      string  true  "ICP filename,example:BE0796544895_202209_01154020.xlsx"
// @Success      200
// @Failure      400
// @Router       /icp/download/{filename} [get]
func DownloadFile(c echo.Context) error {
	tmpDir := viper.GetString("icp.save-dir")
	if !utils.IsDir(tmpDir) {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("The ICP root directory: %s is not exists.", tmpDir))
	}
	filename := c.Param("filename")
	if filename == "" {
		return c.String(http.StatusBadRequest, fmt.Sprintf("The filename must be provided,but was empty."))
	}
	date := strings.Split(filename, "_")[1]
	icpd, err := time.Parse("200601", date)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("File name format not supported(exp: BE0796544895_202209_01154020.xlsx)"))
	}

	icpPath := filepath.Join(tmpDir, icpd.Format("2006"), icpd.Format("01"), filename)

	if !utils.IsExists(icpPath) {
		log.Printf("The icp: %s not found.\n", icpPath)
		return c.String(http.StatusNotFound, fmt.Sprintf("The icp:%s not found.", filename))
	}
	fmt.Println("icpPath", icpPath)

	return c.Attachment(icpPath, filename)
}
