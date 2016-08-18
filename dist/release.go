package dist

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Release represents a published version
type Release struct {
	ID          string
	Version     string
	Description string
	Date        time.Time
}

// DataDir is the path to store files
const DataDir = "./data"

func init() {
	err := os.MkdirAll(DataDir, 0400)
	if err != nil {
		log.Fatal(err)
	}
}

func openReleaseFile() (string, *os.File, error) {
	ts := time.Now().Format("20060102T150405Z0700")
	tc := 0
	for {
		var name string
		if tc == 0 {
			name = ts
		} else {
			name = fmt.Sprintf("%s_%d", ts, tc)
		}

		f, err := os.OpenFile(name+".json", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0400)
		if err != nil {
			if os.IsExist(err) {
				tc = tc + 1
				continue
			}
			return "", nil, err
		}
		return name, f, nil
	}
}

// List list all releases
func List() ([]Release, error) {
	files, err := ioutil.ReadDir(DataDir)
	if err != nil {
		return nil, err
	}

	list := []Release{}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			f, err := os.Open(file.Name())
			if err != nil {
				return nil, err
			}
			item := Release{}
			err = json.NewDecoder(f).Decode(&item)
			if err != nil {
				return nil, fmt.Errorf("unable to decode file: %s", file.Name())
			}
			list = append(list, item)
		}
	}

	return list, nil
}

// Publish uploads & publishes a new version
func Publish(release Release, r io.Reader) (*Release, error) {
	id, f, err := openReleaseFile()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	saved := &Release{
		ID:          id,
		Version:     release.Version,
		Description: release.Description,
		Date:        time.Now(),
	}

	dataf, err := os.OpenFile(DataDir+"/"+id+".dat", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 400)
	if err != nil {
		os.RemoveAll(f.Name())
		return nil, err
	}
	defer dataf.Close()

	_, err = io.Copy(dataf, r)
	if err != nil {
		os.RemoveAll(f.Name())
		return nil, err
	}

	err = json.NewEncoder(f).Encode(saved)
	if err != nil {
		return nil, err
	}

	return saved, nil
}

// Unpublish deletes a published version
func Unpublish(id string) error {
	return nil
}
