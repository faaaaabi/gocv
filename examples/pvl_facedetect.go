// What it does:
//
// This example uses the Intel CV SDK PVL FaceDetect class to detect faces,
// and draw a rectangle around each of them, before displaying them within a Window.
//
// How to run:
//
// pvl_facedetect [camera ID]
//
// 		go run ./examples/pvl_facedetect.go 0
//
// +build example

package main

import (
	"fmt"
	"os"
	"strconv"

	opencv3 ".."
	pvl "../pvl"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("How to run:\n\tpvl_facedetect [camera ID]")
		return
	}

	// parse args
	deviceID, _ := strconv.Atoi(os.Args[1])

	// open webcam
	webcam := opencv3.NewVideoCapture()
	defer webcam.Close()

	if ok := webcam.OpenDevice(deviceID); !ok {
		fmt.Printf("error opening device: %v\n", deviceID)
		return
	}

	// open display window
	window := opencv3.NewWindow("PVL Face Detect")
	defer window.Close()
	
	// prepare input image matrix
	img := opencv3.NewMat()
	defer img.Close()

	// prepare grayscale image matrix
	imgGray := opencv3.NewMat()
	defer imgGray.Close()
	
	// color to draw the rect for detected faces
	blue := opencv3.NewScalar(255, 0, 0, 0)
	green := opencv3.NewScalar(0, 255, 0, 0)

	// load PVL FaceDetector to recognize faces
	fd := pvl.NewFaceDetector()
	defer fd.Close()

	// enable tracking mode for more efficient tracking of video source
	fd.SetTrackingModeEnabled(true)

	fmt.Printf("start reading camera device: %v\n", deviceID)
	for {
		if ok := webcam.Read(img); !ok {
			fmt.Printf("cannot read device %d\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		// convert image to grayscale for detection
		opencv3.CvtColor(img, imgGray, opencv3.ColorBGR2GRAY);
	
		// detect faces
		faces := fd.DetectFaceRect(imgGray)
		fmt.Printf("found %d faces\n", len(faces))

		// draw a rectangle around each face on the original image,
		// along with text identifing as "Human"
		for _, face := range faces {
			// detect smile
			fd.DetectEye(imgGray, face)
			fd.DetectSmile(imgGray, face)

			// set the color of the box based on if the human is smiling
			color := blue
			if face.IsSmiling() {
				color = green
			}

			opencv3.Rectangle(img, face.Rect(), color)

			size := opencv3.GetTextSize("Human", opencv3.FontHersheyPlain, 1.2, 2)
			pt := opencv3.Point{
				X: face.Rect().X + (face.Rect().Width / 2) - (size.Width / 2),
				Y: face.Rect().Y - 2,
			}
			opencv3.PutText(img, "Human", pt, opencv3.FontHersheyPlain, 1.2, color, 2)
		}

		// show the image in the window, and wait 1 millisecond
		window.IMShow(img)
		opencv3.WaitKey(1)
	}
}
