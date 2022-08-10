package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrzack99s/coco-captive-portal/config"
	"github.com/mrzack99s/coco-captive-portal/firewall"
	"github.com/mrzack99s/coco-captive-portal/services"
	"github.com/mrzack99s/coco-captive-portal/types"
	"github.com/mrzack99s/coco-captive-portal/utils"
)

// Headers godoc
// @Summary Get config
// @Description Get config
// @ID get-config
// @Accept   json
// @Tags	Operator
// @Produce  json
// @security ApiKeyAuth
// @Success 200 {object} types.ExtendConfigType
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/config [get]
func (ctl *operatorController) getConfig(c *gin.Context) {
	conf := types.ExtendConfigType{}
	conf.ConfigType = config.Config

	secureIp, _ := utils.GetSecureInterfaceIpv4Addr()
	conf.Status.SecureIPAddress = secureIp
	egressIp, _ := utils.GetEgressInterfaceIpv4Addr()
	conf.Status.EgressIPAddress = egressIp

	c.JSON(http.StatusOK, conf)
}

// Headers godoc
// @Summary Set config
// @Description Set config
// @ID set-config
// @Accept   json
// @Tags	Operator
// @Produce  json
// @security ApiKeyAuth
// @Param params body types.ConfigType true "Parameters"
// @Success 200 {string} string "updated"
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/config [put]
func (ctl *operatorController) setConfig(c *gin.Context) {
	var configs types.ConfigType
	if err := c.ShouldBind(&configs); err != nil {
		msg := "bind a config interface failed"
		config.AppLog.Error().Msg(msg)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": msg,
		})
		return
	}

	tempAllowEndpoint := config.Config.AllowEndpoints
	tempFQDNBlocklist := config.Config.FQDNBlocklist
	tempBypassNetworks := config.Config.BypassedNetworks

	if len(tempAllowEndpoint) > len(configs.AllowEndpoints) {
		diffArray := []types.EndpointType{}
		for _, e1 := range tempAllowEndpoint {
			foundDiff := true
			for _, e2 := range configs.AllowEndpoints {
				if e1 == e2 {
					foundDiff = false
					break
				}
			}

			if foundDiff {
				diffArray = append(diffArray, e1)
			}
		}

		for _, e := range diffArray {
			firewall.DelAllowEndpoint(&e)
		}

	} else {

		diffArray := []types.EndpointType{}
		for _, e1 := range configs.AllowEndpoints {
			foundDiff := true
			for _, e2 := range tempAllowEndpoint {
				if e1 == e2 {
					foundDiff = false
					break
				}
			}

			if foundDiff {
				diffArray = append(diffArray, e1)
			}
		}

		for _, e := range diffArray {
			firewall.AddAllowEndpoint(&e)
		}
	}

	if len(tempFQDNBlocklist) > len(configs.FQDNBlocklist) {
		diffArray := []string{}
		for _, e1 := range tempFQDNBlocklist {
			foundDiff := true
			for _, e2 := range configs.FQDNBlocklist {
				if e1 == e2 {
					foundDiff = false
					break
				}
			}

			if foundDiff {
				diffArray = append(diffArray, e1)
			}
		}

		for _, e := range diffArray {
			firewall.DelFQDNBlacklist(e)
		}

	} else {

		diffArray := []string{}
		for _, e1 := range configs.FQDNBlocklist {
			foundDiff := true
			for _, e2 := range tempFQDNBlocklist {
				if e1 == e2 {
					foundDiff = false
					break
				}
			}

			if foundDiff {
				diffArray = append(diffArray, e1)
			}
		}

		for _, e := range diffArray {
			firewall.AddFQDNBlacklist(e)
		}
	}

	if len(tempBypassNetworks) > len(configs.BypassedNetworks) {
		diffArray := []string{}
		for _, e1 := range tempBypassNetworks {
			foundDiff := true
			for _, e2 := range configs.BypassedNetworks {
				if e1 == e2 {
					foundDiff = false
					break
				}
			}

			if foundDiff {
				diffArray = append(diffArray, e1)
			}
		}

		fmt.Println(diffArray)

		for _, e := range diffArray {
			fmt.Println(e)
			firewall.UnallowAccessBypass(&types.SessionType{
				IPAddress: e,
			})
		}

	} else {

		diffArray := []string{}
		for _, e1 := range configs.BypassedNetworks {
			foundDiff := true
			for _, e2 := range tempBypassNetworks {
				if e1 == e2 {
					foundDiff = false
					break
				}
			}

			if foundDiff {
				diffArray = append(diffArray, e1)
			}
		}

		for _, e := range diffArray {
			firewall.AllowAccessBypass(&types.SessionType{
				IPAddress: e,
			})
		}
	}

	config.Config = configs
	config.UpdateConfig()

	c.String(http.StatusOK, "updated")
}

// Headers godoc
// @Summary Set config with restart system
// @Description Set config with restart system
// @ID set-config-with-restart-system
// @Accept   json
// @Tags	Operator
// @Produce  json
// @security ApiKeyAuth
// @Param params body types.ConfigType true "Parameters"
// @Success 200 {string} string "updated"
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /api/config-with-restart-system [put]
func (ctl *operatorController) setConfigWithRestartSystem(c *gin.Context) {
	var configs types.ConfigType
	if err := c.ShouldBind(&configs); err != nil {
		msg := "bind a config interface failed"
		config.AppLog.Error().Msg(msg)
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": msg,
		})
		return
	}
	config.Config = configs
	config.UpdateConfig()
	services.RestartSystem()

	c.String(http.StatusOK, "updated")
}
