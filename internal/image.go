/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package internal

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/chai2010/webp"
	"github.com/gabriel-vasile/mimetype"
	"github.com/nfnt/resize"
)

// Error que indica que el tipo MIME no es válido
var ErrImageFormatError error = errors.New("invalid image MIME type")

// Recibe una imagen como bytes, detecta el tipo y convierte a image.Image con el decodificador adecuado
func imagenDesdeMime(archivo []byte) (image.Image, error) {
	ctype := mimetype.Detect(archivo)

	switch {
	case ctype.Is("image/png"):
		return png.Decode(bytes.NewReader(archivo))
	case ctype.Is("image/jpeg"):
		return jpeg.Decode(bytes.NewReader(archivo))
	case ctype.Is("image/webp"):
		return webp.Decode(bytes.NewReader(archivo))
	default:
		return nil, ErrImageFormatError
	}

}

func generateImagesFromImage(photo io.Reader) (portadaJpg bytes.Buffer, portadaWebp bytes.Buffer, err error) {
	var portada image.Image
	photoBytes, err := io.ReadAll(photo)
	if err != nil {
		return
	}
	portadaJpg = bytes.Buffer{}
	portadaWebp = bytes.Buffer{}

	ctype := mimetype.Detect(photoBytes)
	if err != nil {
		return
	}

	switch {
	case ctype.Is("image/png"):
		portada, err = png.Decode(bytes.NewReader(photoBytes))
	case ctype.Is("image/jpeg"):
		portada, err = jpeg.Decode(bytes.NewReader(photoBytes))
	case ctype.Is("image/webp"):
		portada, err = webp.Decode(bytes.NewReader(photoBytes))
	default:
		err = ErrImageFormatError
		return
	}

	if err != nil {
		return
	}

	// Resize to 800x450
	portada = resize.Resize(800, 450, portada, resize.Bicubic)

	// Encode as formats
	err = jpeg.Encode(&portadaJpg, portada, &jpeg.Options{Quality: 95})
	if err != nil {
		return
	}
	err = webp.Encode(&portadaWebp, portada, &webp.Options{Quality: 98})
	if err != nil {
		return
	}

	return
}
