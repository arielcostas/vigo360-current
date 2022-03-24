package admin

import (
	"bytes"
	"database/sql"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"strings"

	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
	"golang.org/x/crypto/bcrypt"
)

func ValidatePassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		return true
	}

	if errors.Is(err, bcrypt.ErrHashTooShort) {
		logger.Notice("[validatepassword]: unable to verify password: hash is too short")
	}

	return false
}

func generateImagesFromImage(photo io.Reader, mime *multipart.FileHeader) (portadaJpg bytes.Buffer, portadaWebp bytes.Buffer) {
	var portada image.Image
	var err error
	if strings.HasSuffix(mime.Filename, "png") {
		portada, err = png.Decode(photo)
	} else if strings.HasSuffix(mime.Filename, "jpg") {
		portada, err = jpeg.Decode(photo)
	} else {
		portada, _, err = image.Decode(photo)
	}
	if err != nil {
		log.Fatalln("error imagen 78:" + err.Error())
	}

	// Resize to 800x450
	portada = resize.Resize(800, 450, portada, resize.Bicubic)

	// Encode as formats
	jpeg.Encode(&portadaJpg, portada, &jpeg.Options{Quality: 95})
	webp.Encode(&portadaWebp, portada, &webp.Options{Quality: 98})

	return
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}
