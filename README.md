# Go CKAN
Go wrapper for the [CKAN API](https://docs.ckan.org/en/2.10/contents.html).   
Handles HTTP requests to the CKAN API, provides a structured response.  

https://docs.ckan.org/en/2.10/api/

## About CKAN
The Comprehensive Knowledge Archive Network (CKAN) is an open-source open data portal for the storage and distribution of open data. Initially inspired by the package management capabilities of Debian Linux, CKAN has developed into a powerful data catalogue system that is mainly used by public institutions seeking to share their data with the general public. [Wikipedia](https://en.wikipedia.org/wiki/CKAN).

Some of the organizations using CKAN include:   
- US Government (Data.gov)   
- NHS UK   
- Canadian Government   

## Usage 
### Installation
```bash
go get github.com/ztkent/go-ckan
```

### Example
```go
// Create a new go-ckan client
ckanClient := ckan.NewClient("https://catalog.data.gov/api/3")

// Get the list of recently modified datasets
packageList, err := ckanClient.GetPackageList("metadata_modified desc", 10)

// Get a list of recently created dataset IDs
packageIDList, err := ckanClient.GetPackageIDList("metadata_created desc", 10)

// Search for datasets with the tag "local"
res, err := ckanClient.SearchPackages(
	map[string]string{
		"tags": "local",
	}, "views_recent desc", 10)

// Search for resources with the format "CSV"
res, err := ckanClient.SearchResources(
	map[string]string{
		"format": "CSV",
	}, "views_recent desc", 10)
```

### Endpoints
| Method | Parameters | Return Type | Description |
| --- | --- | --- | --- |
| `SearchPackages` | `query map[string]string`, `sort string`, `limit int` | `PackageSearchResponse`, `error` | Searches for packages based on their metadata. |
| `SearchResources` | `query map[string]string`, `sort string`, `limit int` | `ResourceSearchResponse`, `error` | Searches for resources based on their metadata. |
| `RecentlyChangedPackagesActivityList` | None | `RecentActivityResponse`, `error` | Retrieves an activity stream of recently changed datasets on a site. |
| `GetPackageList` | `sort string`, `limit int` | `PackageIDResponse`, `error` | Retrieves a list of all packages. |
| `GetPackageIDList` | `sort string`, `limit int` | `[]string`, `error` | Retrieves a list of all package IDs. |
| `GetPackageMetadata` | `packageID string` | `PackageMetadataResponse`, `error` | Retrieves metadata for a specific package. |
| `GetResourceID` | `packageID string` | `string`, `error` | Retrieves the resource ID for a specific package. |
| `GetResourceMetadata` | `resourceID string` | `ResourceMetadataResponse`, `error` | Retrieves metadata for a specific resource. |
| `GetGroupList` | `sort string`, `limit int` | `GroupListResponse`, `error` | Retrieves a list of all groups. |
| `GetGroupMetadata` | `groupID string` | `GroupMetadataResponse`, `error` | Retrieves metadata for a specific group. |
| `GetTagList` | `sort string`, `limit int` | `TagListResponse`, `error` | Retrieves a list of all tags. |
| `GetTagMetadata` | `tagName string` | `TagMetadataResponse`, `error` | Retrieves metadata for a specific tag. |
