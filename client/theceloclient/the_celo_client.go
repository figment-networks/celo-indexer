package theceloclient

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	base "github.com/figment-networks/celo-indexer/client"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

const (
	NameTheCelo = "the_celo_client"
)

var (
	_ Client = (*client)(nil)
)

type Client interface {
	base.Client

	GetAllProposals() (*Proposals, error)
	GetProposalByProposalId(string) (*ProposalDetails, error)
}

type client struct {
	baseUrl string
	cc      *http.Client
}

func New(url string) (*client, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify : true},
	}

	cc := &http.Client{Transport: tr}

	return &client{
		baseUrl: url,
		cc:      cc,
	}, nil
}

func (l *client) GetName() string {
	return NameTheCelo
}

func (l *client) Close() {}

func (l *client) GetAllProposals() (*Proposals, error) {
	resp, err := l.cc.Get(fmt.Sprintf("%s?method=proposalList", l.baseUrl))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	proposals := &Proposals{}

	if err = json.Unmarshal(body, proposals); err != nil {
		return nil, err
	}

	return proposals, nil
}

func (l *client) GetProposalByProposalId(proposalId string) (*ProposalDetails, error) {
	proposals, err := l.GetAllProposals()
	if err != nil {
		return nil, err
	}

	proposal, ok := proposals.Items[proposalId]
	if ok {
		return &proposal, nil
	} else {
		return nil, errors.New(fmt.Sprintf("proposal with Id=%s could not be found", proposalId))
	}
}
