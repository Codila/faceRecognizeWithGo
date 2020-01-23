package main

import (
	"log"
	 "image"
	 "image/color"
	"gocv.io/x/gocv"
	"github.com/machinebox/sdk-go/facebox"
	"github.com/hegedustibor/htgo-tts"
//	"time"
//	"fmt"
	"bytes"
)

var (
	faceAlgorithm = "haarcascade_frontalface_default.xml"
	blue = color.RGBA{0, 0, 255, 0}
	fbox = facebox.New("http://localhost:8080")
)
func main() { 
    webcam, err := gocv.VideoCaptureDevice(0)
	// Tuning parameters for faster video processing
	webcam.Set(gocv.VideoCaptureBufferSize, 10)
	//webcam.Set(gocv.VideoCaptureFPS, 50)
	//webcam.Set(gocv.VideoCaptureTriggerDelay, 0)  
	
	
	if err != nil{
			log.Fatalf("unable to init webcam: %v", err)
		}
	
		defer webcam.Close()	

		img := gocv.NewMat()
		defer img.Close()
	
		window := gocv.NewWindow("Face detection in Go")
		defer window.Close()
		
		classifier := gocv.NewCascadeClassifier()
		classifier.Load(faceAlgorithm)
		defer classifier.Close()

		for { 
			 if ok := webcam.Read(&img); !ok || img.Empty() {
				log.Print("unable to read from webcam")
				continue
			}

			rects := classifier.DetectMultiScale(img)
			speech := htgotts.Speech{Folder: "audio", Language: "en"}
					
			for _, r := range rects{
					imgFace := img.Region(r)
					//gocv.IMWrite(fmt.Sprintf("%d.jpg", time.Now().UnixNano()), imgFace)
					buf, err  := gocv.IMEncode(".jpg", imgFace )
					if err != nil{
							log.Printf("unable to encode face img: %v", err)
							continue
					}
					imgFace.Close()
					
					faces, err := fbox.Check(bytes.NewReader(buf)) // Check if the facebox has been pre trained for the encoded Face
					if err != nil {
							log.Printf("unable to check face : %v", err)
							continue
					}
					text := ""
					
					if len(faces) > 0 {
							text = "Hello "+faces[0].Name 		// Get the facial name at initial index from pre-trained facebox statefiles			
					        speech.Speak(text + "Welcome to Digityser") 
					} else {
					speech.Speak("Sorry i dont know you please mention your name") // Greeting if the face not recognized
					}
					// Image border text display parameters tuning										
					size := gocv.GetTextSize(text , gocv.FontHersheyPlain, 3, 2)
					pt := image.Pt(r.Min.X + (r.Min.X / 2) - (size.X / 2), r.Min.Y -2)
					gocv.PutText(&img, text, pt, gocv.FontHersheyPlain, 3, blue, 2)
					gocv.Rectangle(&img, r, blue, 2)
			} 
			window.IMShow(img)
			window.WaitKey(5000)
	    }
}