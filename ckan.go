package ckan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

/*
	Go wrapper for the CKAN API
	https://docs.ckan.org/en/2.10/api/
*/

type CKANAPI interface {
	SearchPackages(query map[string]string, sort string, limit int) (PackageSearchResponse, error)
	SearchResources(query map[string]string, sort string, limit int) (ResourceSearchResponse, error)
	RecentlyChangedPackagesActivityList() (RecentActivityResponse, error)

	GetPackageList(sort string, limit int) (PackageIDResponse, error)
	GetPackageIDList(sort string, limit int) ([]string, error)
	GetPackageMetadata(packageID string) (PackageMetadataResponse, error)

	GetResourceID(packageID string) (string, error)
	GetResourceMetadata(resourceID string) (ResourceMetadataResponse, error)

	GetGroupList(sort string, limit int) (GroupListResponse, error)
	GetGroupMetadata(groupID string) (GroupMetadataResponse, error)

	GetTagList(sort string, limit int) (TagListResponse, error)
	GetTagMetadata(tagName string) (TagMetadataResponse, error)
}

type Client struct {
	BaseURL string
}

type PackageIDResponse struct {
	Help    string `json:"help"`
	Success bool   `json:"success"`
	Result  struct {
		Count        float64                `json:"count"`
		Sort         string                 `json:"sort"`
		Facets       map[string]interface{} `json:"facets"`
		SearchFacets map[string]interface{} `json:"search_facets"`
		Results      []Package              `json:"results"`
	} `json:"result"`
}

type PackageMetadataResponse struct {
	Help    string  `json:"help"`
	Success bool    `json:"success"`
	Result  Package `json:"result"`
}

type ResourceMetadataResponse struct {
	Help    string   `json:"help"`
	Success bool     `json:"success"`
	Result  Resource `json:"result"`
}

type GroupListResponse struct {
	Help    string   `json:"help"`
	Success bool     `json:"success"`
	Result  []string `json:"result"`
}

type GroupMetadataResponse struct {
	Help    string `json:"help"`
	Success bool   `json:"success"`
	Result  Group  `json:"result"`
}

type TagListResponse struct {
	Help    string   `json:"help"`
	Success bool     `json:"success"`
	Result  []string `json:"result"`
}

type TagMetadataResponse struct {
	Help    string `json:"help"`
	Success bool   `json:"success"`
	Result  Tag    `json:"result"`
}

type PackageSearchResponse struct {
	Help    string `json:"help"`
	Success bool   `json:"success"`
	Result  struct {
		Count        float64                `json:"count"`
		Sort         string                 `json:"sort"`
		Facets       map[string]interface{} `json:"facets"`
		SearchFacets map[string]interface{} `json:"search_facets"`
		Results      []Package              `json:"results"`
	} `json:"result"`
}

type ResourceSearchResponse struct {
	Help    string `json:"help"`
	Success bool   `json:"success"`
	Result  struct {
		Count        float64                `json:"count"`
		Sort         string                 `json:"sort"`
		Facets       map[string]interface{} `json:"facets"`
		SearchFacets map[string]interface{} `json:"search_facets"`
		Results      []Package              `json:"results"`
	} `json:"result"`
}

type RecentActivityResponse struct {
	Help    string     `json:"help"`
	Success bool       `json:"success"`
	Result  []Activity `json:"result"`
}

type Package struct {
	Maintainer             string        `json:"maintainer"`
	OwnerOrg               string        `json:"owner_org"`
	Groups                 []interface{} `json:"groups"`
	CreatorUserID          string        `json:"creator_user_id"`
	MetadataModified       string        `json:"metadata_modified"`
	State                  string        `json:"state"`
	Tags                   []Tag         `json:"tags"`
	RelationshipsAsSubject []interface{} `json:"relationships_as_subject"`
	Name                   string        `json:"name"`
	Private                bool          `json:"private"`
	Version                string        `json:"version"`
	LicenseID              string        `json:"license_id"`
	LicenseTitle           string        `json:"license_title"`
	MetadataCreated        string        `json:"metadata_created"`
	NumTags                int           `json:"num_tags"`
	Organization           Organization  `json:"organization"`
	Type                   string        `json:"type"`
	ID                     string        `json:"id"`
	IsOpen                 bool          `json:"isopen"`
	Notes                  string        `json:"notes"`
	Title                  string        `json:"title"`
	Resources              []Resource    `json:"resources"`
	RelationshipsAsObject  []interface{} `json:"relationships_as_object"`
	Author                 string        `json:"author"`
	AuthorEmail            string        `json:"author_email"`
	MaintainerEmail        string        `json:"maintainer_email"`
	NumResources           int           `json:"num_resources"`
	URL                    string        `json:"url"`
	Extras                 []Extra       `json:"extras"`
}

type Resource struct {
	CacheURL         string `json:"cache_url"`
	Format           string `json:"format"`
	MetadataModified string `json:"metadata_modified"`
	MimetypeInner    string `json:"mimetype_inner"`
	NoRealName       bool   `json:"no_real_name"`
	Position         int    `json:"position"`
	Size             int    `json:"size"`
	State            string `json:"state"`
	URL              string `json:"url"`
	Created          string `json:"created"`
	Mimetype         string `json:"mimetype"`
	CacheLastUpdated string `json:"cache_last_updated"`
	Hash             string `json:"hash"`
	ResourceType     string `json:"resource_type"`
	Description      string `json:"description"`
	ID               string `json:"id"`
	LastModified     string `json:"last_modified"`
	Name             string `json:"name"`
	PackageID        string `json:"package_id"`
	URLType          string `json:"url_type"`
	DescribedBy      string `json:"describedBy,omitempty"`
	DescribedByType  string `json:"describedByType,omitempty"`
}

type Organization struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	IsOrganization bool   `json:"is_organization"`
	ApprovalStatus string `json:"approval_status"`
	State          string `json:"state"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	Description    string `json:"description"`
	ImageURL       string `json:"image_url"`
	Created        string `json:"created"`
}

type Group struct {
	Description string `json:"description"`
	Created     string `json:"created"`
	Title       string `json:"title"`
	Name        string `json:"name"`
	IsOrg       bool   `json:"is_organization"`
	State       string `json:"state"`
	ImageURL    string `json:"image_url"`
	Type        string `json:"type"`
	ID          string `json:"id"`
}

type Activity struct {
	ActivityID   string                 `json:"id"`
	Timestamp    string                 `json:"timestamp"`
	UserID       string                 `json:"user_id"`
	ObjectID     string                 `json:"object_id"`
	RevisionID   string                 `json:"revision_id"`
	ActivityType string                 `json:"activity_type"`
	Data         map[string]interface{} `json:"data"`
}

type Tag struct {
	Name string `json:"name"`
}

type Extra struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func NewClient(baseURL string) Client {
	return Client{BaseURL: baseURL}
}

// Search for packages based on their metadata
func (ckan Client) SearchPackages(query map[string]string, sort string, limit int) (PackageSearchResponse, error) {
	params := setSearchParams(query, sort, limit)
	resp, err := http.Get(fmt.Sprintf("%s/action/package_search?%s", ckan.BaseURL, params.Encode()))
	if err != nil {
		return PackageSearchResponse{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PackageSearchResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return PackageSearchResponse{}, fmt.Errorf("Failed to search packages: %s", body)
	}

	var res PackageSearchResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return PackageSearchResponse{}, err
	}
	if !res.Success {
		return PackageSearchResponse{}, fmt.Errorf("Failed to search packages: %s", body)
	}

	return res, nil
}

// Search for resources based on their metadata
func (ckan Client) SearchResources(query map[string]string, sort string, limit int) (ResourceSearchResponse, error) {
	params := setSearchParams(query, sort, limit)
	resp, err := http.Get(fmt.Sprintf("%s/action/resource_search?%s", ckan.BaseURL, params.Encode()))
	if err != nil {
		return ResourceSearchResponse{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ResourceSearchResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return ResourceSearchResponse{}, fmt.Errorf("Failed to search resources: %s", body)
	}

	var res ResourceSearchResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return ResourceSearchResponse{}, err
	}
	if !res.Success {
		return ResourceSearchResponse{}, fmt.Errorf("Failed to search resources: %s", body)
	}

	return res, nil
}

// Get a set of recently changed datasets
func (ckan Client) RecentlyChangedPackagesActivityList() (RecentActivityResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api/3/action/recently_changed_packages_activity_list", ckan.BaseURL))
	if err != nil {
		return RecentActivityResponse{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return RecentActivityResponse{}, fmt.Errorf("Failed to get recently changed packages activity list: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return RecentActivityResponse{}, err
	}

	var res RecentActivityResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return RecentActivityResponse{}, err
	}
	if !res.Success {
		return RecentActivityResponse{}, fmt.Errorf("Failed to get recently changed packages activity list: %s", body)
	}

	return res, nil
}

// Get a list of dataset IDs
func (ckan Client) GetPackageIDList(sort string, limit int) ([]string, error) {
	res, err := ckan.GetPackageList(sort, limit)
	if err != nil {
		return nil, err
	}
	packages := make([]string, 0)
	for _, r := range res.Result.Results {
		packages = append(packages, r.ID)
	}

	return packages, nil
}

// Get a list of datasets
// Example Sort: "views_recent desc", "metadata_modified desc"
func (ckan Client) GetPackageList(sort string, limit int) (PackageIDResponse, error) {
	params := setSearchParams(nil, sort, limit)
	resp, err := http.Get(fmt.Sprintf("%s/action/package_list?%s", ckan.BaseURL, params.Encode()))
	if err != nil {
		return PackageIDResponse{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PackageIDResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return PackageIDResponse{}, fmt.Errorf("Failed to get package list: %s", body)
	}

	res := PackageIDResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return PackageIDResponse{}, err
	}
	if !res.Success {
		return PackageIDResponse{}, fmt.Errorf("Failed to get package list: %s", body)
	}

	return res, nil
}

// Get metadata about a dataset ID
func (ckan Client) GetPackageMetadata(packageID string) (PackageMetadataResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s/action/package_show?id=%s", ckan.BaseURL, packageID))
	if err != nil {
		return PackageMetadataResponse{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return PackageMetadataResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return PackageMetadataResponse{}, fmt.Errorf("Failed to get package metadata: %s", body)
	}

	var res PackageMetadataResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return PackageMetadataResponse{}, err
	}
	if !res.Success {
		return PackageMetadataResponse{}, fmt.Errorf("Failed to get package metadata: %s", body)
	}

	return res, nil
}

// Get the resource ID for a dataset
func (ckan Client) GetResourceID(packageID string) (string, error) {
	packageMetadata, err := ckan.GetPackageMetadata(packageID)
	if err != nil {
		return "", err
	}
	resources := packageMetadata.Result.Resources
	if len(resources) == 0 {
		return "", fmt.Errorf("No resources found for package %s", packageID)
	}
	for _, resource := range resources {
		if resource.ID != "" {
			return resource.ID, nil
		}
	}

	return "", nil
}

// Get metadata about a resource ID
func (ckan Client) GetResourceMetadata(resourceID string) (ResourceMetadataResponse, error) {
	data := map[string]string{"id": resourceID}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return ResourceMetadataResponse{}, err
	}
	resp, err := http.Post(fmt.Sprintf("%s/action/resource_show", ckan.BaseURL), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return ResourceMetadataResponse{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ResourceMetadataResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return ResourceMetadataResponse{}, fmt.Errorf("Failed to get resource metadata: %s", body)
	}

	var res ResourceMetadataResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return ResourceMetadataResponse{}, err
	}
	if !res.Success {
		return ResourceMetadataResponse{}, fmt.Errorf("Failed to get resource metadata: %s", body)
	}

	return res, nil
}

// Get a list of available groups
func (ckan Client) GetGroupList(sort string, limit int) (GroupListResponse, error) {
	params := setSearchParams(nil, sort, limit)
	resp, err := http.Get(fmt.Sprintf("%s/action/group_list?%s", ckan.BaseURL, params.Encode()))
	if err != nil {
		return GroupListResponse{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return GroupListResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return GroupListResponse{}, fmt.Errorf("Failed to get group list: %s", body)
	}

	res := GroupListResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return GroupListResponse{}, err
	}
	if !res.Success {
		return GroupListResponse{}, fmt.Errorf("Failed to get group list: %s", body)
	}

	return res, nil
}

// Get metadata about a group ID
func (ckan Client) GetGroupMetadata(groupID string) (GroupMetadataResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s/action/group_show?id=%s", ckan.BaseURL, groupID))
	if err != nil {
		return GroupMetadataResponse{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return GroupMetadataResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return GroupMetadataResponse{}, fmt.Errorf("Failed to get group metadata: %s", body)
	}

	var res GroupMetadataResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return GroupMetadataResponse{}, err
	}
	if !res.Success {
		return GroupMetadataResponse{}, fmt.Errorf("Failed to get group metadata: %s", body)
	}

	return res, nil
}

// Get a list of available tags
func (ckan Client) GetTagList(sort string, limit int) (TagListResponse, error) {
	params := setSearchParams(nil, sort, limit)
	resp, err := http.Get(fmt.Sprintf("%s/action/tag_list?%s", ckan.BaseURL, params.Encode()))
	if err != nil {
		return TagListResponse{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TagListResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return TagListResponse{}, fmt.Errorf("Failed to get tag list: %s", body)
	}

	res := TagListResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return TagListResponse{}, err
	}
	if !res.Success {
		return TagListResponse{}, fmt.Errorf("Failed to get tag list: %s", body)
	}

	return res, nil
}

// Get metadata about a tag
func (ckan Client) GetTagMetadata(tagName string) (TagMetadataResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s/action/tag_show?id=%s", ckan.BaseURL, tagName))
	if err != nil {
		return TagMetadataResponse{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TagMetadataResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return TagMetadataResponse{}, fmt.Errorf("Failed to get tag metadata: %s", body)
	}

	var res TagMetadataResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return TagMetadataResponse{}, err
	}
	if !res.Success {
		return TagMetadataResponse{}, fmt.Errorf("Failed to get tag metadata: %s", body)
	}

	return res, nil
}

func setSearchParams(query map[string]string, sort string, limit int) url.Values {
	var searchTerms []string
	for key, value := range query {
		searchTerms = append(searchTerms, fmt.Sprintf("%s:%s", key, value))
	}
	searchQuery := strings.Join(searchTerms, "+")

	params := url.Values{}
	if searchQuery != "" {
		params.Add("q", searchQuery)
	}
	if sort != "" {
		params.Add("sort", sort)
	}
	if limit != 0 {
		params.Add("rows", strconv.Itoa(limit))
	}
	return params
}
