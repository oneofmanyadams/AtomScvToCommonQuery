package main

import (	
	"fmt"
	"strconv"
	"strings"

	"io/ioutil"
	"path/filepath"

	"data/atomsvc"
	"data/commonquery"
)

func main() {
	// Determine all files in active directory
	all_files, read_dir_err := ioutil.ReadDir("./")
	if read_dir_err != nil {
		fmt.Println(read_dir_err)
	}

	// Loop through files in active directory
	for _, a_file := range all_files {
		
		// Get file extension of file
		file_extension := filepath.Ext(a_file.Name())
		
		// We only care about atomsvc files
		if file_extension == ".atomsvc" {
		
			// Initialize one MultiQuery per .atomsvc file
			queries := commonquery.NewMultiQuery()
			
			// Load atomsvc data
			atomsvc_data := atomsvc.FromFile(a_file.Name())
			
			// Loop through each Collection record in the atomsvc data
			for collection_count, collection := range atomsvc_data.AllCollections() {
		
				// Create new CommonQuery for each atomsvc collection
				new_query := commonquery.NewCommonQuery()
				
				// Directly copy atomsvc data HostPath to CommonQuery HostPath
				new_query.HostPath = collection.HostPath
		
				// Copy over atomsvc data options to CommonQuery options.
				for option_name, option_values := range collection.Options {
					for _, val := range option_values {
						new_query.AddOption(option_name, val)
					}
				}
				
				// Have the CommonQuery assemble it's full path from copied over atomsvc data
				new_query.BuildFullPath()

				// Add CommonQuery built from atomsvc-collection-data into MultiQuery
				queries.AddCommonQuery(strconv.Itoa(collection_count), new_query)
			}

			// Save CommonQuery data using atomsvc filename, but:
			//	-Replaces ".atomsvc" extension with ".json" extension
			//	-Appends "QC_" prefix to file name.
			queries.SaveTo("CQ_"+strings.Replace(a_file.Name(), file_extension, ".json", -1))
		}
	}

}