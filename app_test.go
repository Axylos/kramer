package main

import (
	"bytes"
	"log"
	"os"
	"testing"

	"fyne.io/fyne/test"
	"github.com/axylos/kramer_pager/kramer"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	ta := test.NewApp()
	k := kramer.NewKramer(ta)

	k.Run()

	assert.Equal(t, "Kramer Pager", k.Win.Title())
}

func TestButton(t *testing.T) {
	ta := test.NewApp()
	k := kramer.NewKramer(ta)

	println(string(k.Win.Title()))
	k.Run()

	var buf bytes.Buffer
	log.SetOutput(&buf)
	test.Tap(k.Buttons["signup"])
	log.SetOutput(os.Stderr)
	assert.Equal(t, true, k.Tapped)
	//assert.Equal(t, "Kramer Pager", buf.String())
}
