package utils

import (
	"fmt"
	"testing"
)

func TestGetRedirectUrl_basepage(t *testing.T) {
	var redir_url = GetRedirectUrl("https://www.google.com/search?sxsrf=ACYBGNRZr5M9ppUzFZQKPuaYOQtQL5EMtQ%3A1572167872215&source=hp&ei=wGC1XYDPCoWRmwWZh4b4Bg&q=how+to+write+go+unit+tests&oq=how+to+write+go+unit+tests&gs_l=psy-ab.3..0i22i30l2.1869.6449..6577...4.0..0.125.2231.25j5......0....1..gws-wiz.......0i131j0j35i39i19j0i10j0i203j0i22i10i30j0i19j0i13i30i19j0i8i13i30i19j33i22i29i30j0i13i30.jSyRCAg_ZBE&ved=0ahUKEwiAr7HcjbzlAhWFyKYKHZmDAW8Q4dUDCAc&uact=5")
	if redir_url != "/" {
		t.Fail()
	}
}

func TestGetRedirectUrl_maps(t *testing.T) {
	var redir_url = GetRedirectUrl("https://www.google.com/maps/place/New+York+City,+New+York,+USA/@40.6974034,-74.119763,11z/data=!3m1!4b1!4m5!3m4!1s0x89c24fa5d33f083b:0xc80b8f06e177fe62!8m2!3d40.7127753!4d-74.0059728")
	fmt.Println(redir_url)
	if redir_url != "/maps/place/New+York+City,+New+York,+USA/@40.6974034,-74.119763,11z/data=!3m1!4b1!4m5!3m4!1s0x89c24fa5d33f083b:0xc80b8f06e177fe62!8m2!3d40.7127753!4d-74.0059728" {
		t.Fail()
	}
}

func TestGetRedirectUrl_malformedUrl(t *testing.T) {
	var redir_url = GetRedirectUrl("this_is_no_url")
	fmt.Println(redir_url)
	if redir_url != "/" {
		t.Fail()
	}
}
