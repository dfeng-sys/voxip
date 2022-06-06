package voxip_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"voxip.com/voxip/voxip"
)

var (
	mmdbFile = "../db/GeoLite2-Country.mmdb"
)

func TestCheckWhitelist_PositiveScenario(t *testing.T) {

	voxip.Init(mmdbFile)

	req := voxip.Request{
		Ip:        "81.2.69.142",
		Whitelist: []string{"US", "GB", "CN"},
	}
	resp := voxip.CheckWhitelist(req)
	assert.Empty(t, resp.Error)
	assert.True(t, resp.Whitelisted)
	assert.Equal(t, "GB", resp.Country)

	req.Ip = "35.225.41.128"
	resp = voxip.CheckWhitelist(req)
	assert.Empty(t, resp.Error)
	assert.True(t, resp.Whitelisted)
	assert.Equal(t, "US", resp.Country)
}

func TestCheckWhitelist_NegativeScenario(t *testing.T) {

	voxip.Init(mmdbFile)

	req := voxip.Request{
		Ip:        "81.2.69.142",
		Whitelist: []string{"AU", "DE", "JP"},
	}
	resp := voxip.CheckWhitelist(req)
	assert.Empty(t, resp.Error)
	assert.False(t, resp.Whitelisted)
	assert.Equal(t, "GB", resp.Country)

	req.Ip = "35.225.41.128"
	resp = voxip.CheckWhitelist(req)
	assert.Empty(t, resp.Error)
	assert.False(t, resp.Whitelisted)
	assert.Equal(t, "US", resp.Country)
}

func TestCheckWhitelist_InvalidIp(t *testing.T) {

	voxip.Init(mmdbFile)

	req := voxip.Request{
		Ip:        "4.8.15.16.23.42",
		Whitelist: []string{},
	}
	resp := voxip.CheckWhitelist(req)
	assert.Equal(t, "received invalid IP address", resp.Error)
	assert.False(t, resp.Whitelisted)
	assert.Empty(t, resp.Country)
}

func TestCheckWhitelist_EmptyIp(t *testing.T) {

	voxip.Init(mmdbFile)

	req := voxip.Request{
		Ip:        "",
		Whitelist: []string{},
	}
	resp := voxip.CheckWhitelist(req)
	assert.Equal(t, "received null or empty IP address", resp.Error)
	assert.False(t, resp.Whitelisted)
	assert.Empty(t, resp.Country)
}

func TestCheckWhitelist_EmptyWhitelist(t *testing.T) {

	voxip.Init(mmdbFile)

	req := voxip.Request{
		Ip:        "35.225.41.128",
		Whitelist: []string{},
	}
	resp := voxip.CheckWhitelist(req)
	assert.Equal(t, "received null or empty whitelist", resp.Error)
	assert.False(t, resp.Whitelisted)
	assert.Empty(t, resp.Country)

	req.Whitelist = nil
	resp = voxip.CheckWhitelist(req)
	assert.Equal(t, "received null or empty whitelist", resp.Error)
	assert.False(t, resp.Whitelisted)
	assert.Empty(t, resp.Country)
}
