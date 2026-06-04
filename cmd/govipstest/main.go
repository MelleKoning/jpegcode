package main

import (
	"fmt"
	"log"
	"os"

	"github.com/davidbyttow/govips/v2/vips"
)

func main() {
	// 1. Startup the libvips C-engine threadpool
	err := vips.Startup(nil)
	defer vips.Shutdown()
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("--- Processing Lion Head Big (Baseline, 10MB) ---")
	processImage("data/lionheadbig.jpg", "data/output_lionheadbig_to_progressive.jpg", true)

	fmt.Println("\n--- Processing Lion Face (Progressive, 10MB) ---")
	processImage("data/lionheadbig.jpg", "data/output_lionheadbig_nonprogressive.jpg", false)
}

func processImage(inputPath, outputPath string, makeProgressive bool) {
	// 2. Load the file directly into C-allocated memory
	img, err := vips.NewImageFromFile(inputPath)
	if err != nil {
		log.Fatalf("Failed to load image %s: %v", inputPath, err)
	}
	defer img.Close() // CRITICAL: Frees the C memory instantly when this function exits

	// 3. Read metadata properties native to vips
	width := img.Width()
	height := img.Height()
	bands := img.Bands() // 3 means RGB, 4 means RGBA
	fmt.Printf("Loaded: %s | Dimensions: %dx%d | Channels: %d\n", inputPath, width, height, bands)

	// Optional: Let's do a fast native downscale operation just to demonstrate libvips speed
	// This uses a high-quality Lanczos3 kernel under the hood
	/*if width > 2000 {
		fmt.Println("Image is large, resizing down by 50%...")
		err = img.Resize(0.5, vips.KernelLanczos3)
		if err != nil {
			log.Fatalf("Failed to resize: %v", err)
		}
	}*/

	// 4. Configure our granular JPEG encoder parameters
	// This gives us the deep control over the final disk format
	ep := vips.NewJpegExportParams()
	ep.Quality = 95          // 85 is good for web, we use 95
	ep.StripMetadata = true  // Strips EXIF/JFIF header bloat
	ep.OptimizeCoding = true // Computes optimal Huffman tables

	if makeProgressive {
		ep.Interlace = true // This converts the Baseline lionheadbig into a progressive JPEG!
		fmt.Println("Setting export configuration to Progressive layout.")
	} else {
		ep.Interlace = false // Force baseline structure
		fmt.Println("Setting export configuration to Baseline layout.")
	}

	// Force 4:4:4 chroma subsampling if you want zero color bleeding on edges,
	// or leave it default for standard web 4:2:0 subsampling.
	//ep.SubsampleMode = C.VIPS_FOREIGN_SUBSAMPLE_AUTO

	// 5. Execute the export graph into compressed bytes
	jpegBytes, _, err := img.ExportJpeg(ep)
	if err != nil {
		log.Fatalf("Failed to export JPEG: %v", err)
	}

	// 6. Write the raw bytes to disk
	err = os.WriteFile(outputPath, jpegBytes, 0644)
	if err != nil {
		log.Fatalf("Failed to write file to disk: %v", err)
	}

	// Check final output size
	fi, _ := os.Stat(outputPath)
	fmt.Printf("Saved to: %s | New file size: %.2f MB\n", outputPath, float64(fi.Size())/1024/1024)
}
