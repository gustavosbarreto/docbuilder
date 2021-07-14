// +build internal_api

package client

import (
	"fmt"
	"net/http"

	"github.com/shellhub-io/shellhub/pkg/models"
	"go.uber.org/multierr"
)

const (
	apiHost   = "api"
	apiPort   = 8080
	apiScheme = "http"
)

type Client interface {
	commonAPI
	internalAPI
}

type internalAPI interface {
	LookupDevice()
	GetPublicKey(fingerprint, tenant string) (*models.PublicKey, error)
	CreatePrivateKey() (*models.PrivateKey, error)
	EvaluateKey(fingerprint string, dev *models.Device) (bool, error)
	DevicesOffline(id string) error
	FirewallEvaluate(lookup map[string]string) []error
	PatchSessions(uid string) []error
	FinishSession(uid string) []error
	RecordSession(session *models.SessionRecorded, recordURL string)
	Lookup(lookup map[string]string) (string, []error)
	DeviceLookup(lookup map[string]string) (*models.Device, []error)
	GetNamespaceByName(tenant string) (*models.Namespace, error)
}

func (c *client) LookupDevice() {
}

func (c *client) GetPublicKey(fingerprint, tenant string) (*models.PublicKey, error) {
	var pubKey *models.PublicKey
	resp, _, errs := c.http.Get(buildURL(c, fmt.Sprintf("/internal/sshkeys/public-keys/%s/%s", fingerprint, tenant))).EndStruct(&pubKey)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	if resp.StatusCode == 404 {
		return nil, ErrNotFound
	}

	return pubKey, nil
}

func (c *client) EvaluateKey(fingerprint string, dev *models.Device) (bool, error) {
	var evaluate *bool

	resp, _, errs := c.http.Post(buildURL(c, fmt.Sprintf("/internal/sshkeys/public-keys/evaluate/%s", fingerprint))).Send(dev).EndStruct(&evaluate)
	if len(errs) > 0 {
		var err error
		for _, e := range errs {
			err = multierr.Append(err, e)
		}

		return false, err
	}

	if resp.StatusCode == 200 {
		return *evaluate, nil
	}

	return false, nil
}

func (c *client) CreatePrivateKey() (*models.PrivateKey, error) {
	var privKey *models.PrivateKey
	_, _, errs := c.http.Post(buildURL(c, "/internal/sshkeys/private-keys")).EndStruct(&privKey)
	if len(errs) > 0 {
		return nil, errs[0]
	}

	return privKey, nil
}

func (c *client) GetNamespaceByName(tenant string) (*models.Namespace, error) {
	var namespace *models.Namespace
	resp, _, errs := c.http.Get(buildURL(c, fmt.Sprintf("/internal/namespaces/%s", tenant))).EndStruct(&namespace)
	if len(errs) > 0 {
		return nil, ErrConnectionFailed
	}

	if resp.StatusCode == 400 {
		return nil, ErrNotFound
	} else if resp.StatusCode == 200 {
		return namespace, nil
	}

	return nil, ErrUnknown
}

func (c *client) DevicesOffline(id string) error {
	_, _, errs := c.http.Post(buildURL(c, fmt.Sprintf("/internal/devices/%s/offline", id))).End()
	if len(errs) > 0 {
		return errs[0]
	}

	return nil
}

func (c *client) FirewallEvaluate(lookup map[string]string) []error {
	if res, _, errs := c.http.Get(buildURL(c, "/internal/firewall/rules/evaluate")).Query(lookup).End(); res.StatusCode != http.StatusOK {
		return errs
	}

	return nil
}

func (c *client) PatchSessions(uid string) []error {
	_, _, errs := c.http.Patch(buildURL(c, fmt.Sprintf("/internal/sessions/"+uid))).Send(&models.Status{
		Authenticated: true,
	}).End()

	return errs
}

func (c *client) FinishSession(uid string) []error {
	_, _, errs := c.http.Post(buildURL(c, fmt.Sprintf("/internal/sessions/%s/finish", uid))).End()

	return errs
}

func (c *client) RecordSession(session *models.SessionRecorded, recordURL string) {
	c.http.Post(fmt.Sprintf("http://"+recordURL+"/internal/sessions/%s/record", session.UID)).Send(&session).End()
}

func (c *client) Lookup(lookup map[string]string) (string, []error) {
	var device struct {
		UID string `json:"uid"`
	}

	if res, _, errors := c.http.Get(buildURL(c, "/internal/lookup")).Query(lookup).EndStruct(&device); res.StatusCode != http.StatusOK {
		return "", errors
	}

	return device.UID, nil
}

func (c *client) DeviceLookup(lookup map[string]string) (*models.Device, []error) {
	var device *models.Device

	if res, _, errors := c.http.Get(buildURL(c, "/internal/lookup")).Query(lookup).EndStruct(&device); res.StatusCode != http.StatusOK {
		return nil, errors
	}

	return device, nil
}
