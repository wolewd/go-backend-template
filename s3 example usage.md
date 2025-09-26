```go
package main

import (
	"fmt"
	"log"
	"time"
	"your-project/config"
	"your-project/utils"
)

func main() {
	// Initialize S3 client
	s3Client, err := config.ConnectS3()
	if err != nil {
		log.Fatal("Failed to connect to S3:", err)
	}

	// Create S3 helper
	s3Helper := utils.NewS3Helper(s3Client)
	
	bucketName := "my-bucket"
	
	// Example file data (in real app, this comes from uploaded file)
	fileData := []byte("This is example file content for testing")
	
	fmt.Println("=== S3 Helper Examples ===\n")

	// 1. Upload file with original name
	fmt.Println("1. Upload with original filename:")
	fileName1, err := s3Helper.UploadFile(bucketName, "documents/report.pdf", fileData, "application/pdf")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("✅ Uploaded: %s\n", fileName1)
	}
	
	// 2. Upload file with UUID name (NEW METHOD)
	fmt.Println("\n2. Upload with UUID filename:")
	fileName2, err := s3Helper.UploadFileWithUUID(bucketName, "photo.png", fileData, "image/png")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("✅ Uploaded with UUID: %s\n", fileName2)
	}
	
	fileName3, err := s3Helper.UploadFileWithUUID(bucketName, "invoice.pdf", fileData, "application/pdf")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("✅ Uploaded with UUID: %s\n", fileName3)
	}

	// 3. Check if file exists
	fmt.Println("\n3. Check if files exist:")
	exists1, err := s3Helper.FileExists(bucketName, fileName1)
	if err != nil {
		log.Printf("Error checking file: %v", err)
	} else {
		fmt.Printf("File '%s' exists: %t\n", fileName1, exists1)
	}
	
	exists2, err := s3Helper.FileExists(bucketName, fileName2)
	if err != nil {
		log.Printf("Error checking file: %v", err)
	} else {
		fmt.Printf("File '%s' exists: %t\n", fileName2, exists2)
	}

	// 4. Get public download link
	fmt.Println("\n4. Get public download links:")
	publicURL1 := s3Helper.DownloadPublicFile(bucketName, fileName1)
	fmt.Printf("Public URL 1: %s\n", publicURL1)
	
	publicURL2 := s3Helper.DownloadPublicFile(bucketName, fileName2)
	fmt.Printf("Public URL 2: %s\n", publicURL2)

	// 5. Get private download link (with expiry)
	fmt.Println("\n5. Get private download links (expires in 1 hour):")
	privateURL1, err := s3Helper.DownloadPrivateFile(bucketName, fileName1, time.Hour)
	if err != nil {
		log.Printf("Error generating private URL: %v", err)
	} else {
		fmt.Printf("Private URL 1: %s\n", privateURL1)
	}
	
	privateURL2, err := s3Helper.DownloadPrivateFile(bucketName, fileName2, time.Hour)
	if err != nil {
		log.Printf("Error generating private URL: %v", err)
	} else {
		fmt.Printf("Private URL 2: %s\n", privateURL2)
	}

	// 6. Delete files
	fmt.Println("\n6. Delete files:")
	err = s3Helper.DeleteFile(bucketName, fileName1)
	if err != nil {
		log.Printf("Error deleting file 1: %v", err)
	} else {
		fmt.Printf("✅ Deleted: %s\n", fileName1)
	}
	
	err = s3Helper.DeleteFile(bucketName, fileName2)
	if err != nil {
		log.Printf("Error deleting file 2: %v", err)
	} else {
		fmt.Printf("✅ Deleted: %s\n", fileName2)
	}

	// 7. Check if files still exist after deletion
	fmt.Println("\n7. Check files after deletion:")
	exists1After, err := s3Helper.FileExists(bucketName, fileName1)
	if err != nil {
		log.Printf("Error checking file: %v", err)
	} else {
		fmt.Printf("File '%s' exists after deletion: %t\n", fileName1, exists1After)
	}
	
	exists2After, err := s3Helper.FileExists(bucketName, fileName2)
	if err != nil {
		log.Printf("Error checking file: %v", err)
	} else {
		fmt.Printf("File '%s' exists after deletion: %t\n", fileName2, exists2After)
	}
}
```

/*
EXAMPLE OUTPUT:
=== S3 Helper Examples ===

1. Upload with original filename:
✅ Uploaded: documents/report.pdf

2. Upload with UUID filename:
✅ Uploaded with UUID: 01912345-6789-7abc-def0-123456789012.png
✅ Uploaded with UUID: 01912346-678a-7bcd-ef01-234567890123.pdf

3. Check if files exist:
File 'documents/report.pdf' exists: true
File '01912345-6789-7abc-def0-123456789012.png' exists: true

4. Get public download links:
Public URL 1: my-bucket/documents/report.pdf
Public URL 2: my-bucket/01912345-6789-7abc-def0-123456789012.png

5. Get private download links (expires in 1 hour):
Private URL 1: https://your-s3-endpoint.com/my-bucket/documents/report.pdf?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=...
Private URL 2: https://your-s3-endpoint.com/my-bucket/01912345-6789-7abc-def0-123456789012.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=...

6. Delete files:
✅ Deleted: documents/report.pdf
✅ Deleted: 01912345-6789-7abc-def0-123456789012.png

7. Check files after deletion:
File 'documents/report.pdf' exists after deletion: false
File '01912345-6789-7abc-def0-123456789012.png' exists after deletion: false
*/