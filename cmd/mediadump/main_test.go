package main

import (
	"reflect"
	"testing"
)

func Test_extractMediaLinks(t *testing.T) {
	type args struct {
		content string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "simple test",
			args: args{
				content: "<style>#html-body [data-pb-style=CLQALMO]{justify-content:flex-start;display:flex;flex-direction:column;background-position:left top;background-size:cover;background-repeat:no-repeat;background-attachment:scroll}#html-body [data-pb-style=AO6UGUU],#html-body [data-pb-style=NYELC6K]{max-width:100%;height:auto}</style><div data-content-type=\"row\" data-appearance=\"contained\" data-element=\"main\"><div data-enable-parallax=\"0\" data-parallax-speed=\"0.5\" data-background-images=\"{}\" data-background-type=\"image\" data-video-loop=\"true\" data-video-play-only-visible=\"true\" data-video-lazy-load=\"true\" data-video-fallback-src=\"\" data-element=\"inner\" data-pb-style=\"CLQALMO\"><div class=\"c-plp-cms-block\" data-content-type=\"hot_plp_block\" data-appearance=\"default\" data-hf-2col=\"false\" data-element=\"main\"><figure><img class=\"pagebuilder-mobile-hidden\" width=\"300\" height=\"460\" src=\"{{media url=wysiwyg/plp/hf/plp-grid-image-desktop.jpg}}\" alt=\"Girl with hydroflask bottle sitting on car\" data-element=\"image_element_desktop\" data-pb-style=\"AO6UGUU\"><img class=\"pagebuilder-mobile-only\" width=\"710\" height=\"1015\" src=\"{{media url=wysiwyg/plp/hf/plp-grid-image-mobile.jpg}}\" alt=\"Girl with hydroflask bottle sitting on car\" data-element=\"image_element_mobile\" data-pb-style=\"NYELC6K\"></figure><div class=\"c-plp-cms-block__content\"><div class=\"c-plp-cms-block__content--light\" data-element=\"inner\"><p class=\"c-plp-cms-block__heading\" data-element=\"title\"></p><div class=\"c-plp-cms-block__copy\" data-element=\"description\">Alii autem, quibus ego assentior, cum a sapiente delectus</div></div><div class=\"c-plp-cms-block__cta\"><a class=\"pagebuilder-button-customize\" href=\"/\" target=\"\" data-link-type=\"default\" data-element=\"link\"><span class=\"a-btn-customize__container\"><span class=\"a-btn-customize__icon icon-customize\" aria-hidden=\"true\"></span><span class=\"a-btn-customize__label\" data-element=\"link_text\">Customize Yours!</span></span></a></div></div></div></div></div>",
			},
			want: []string{
				"/media/wysiwyg/plp/hf/plp-grid-image-desktop.jpg",
				"/media/wysiwyg/plp/hf/plp-grid-image-mobile.jpg",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractMediaLinks(tt.args.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractMediaLinks() = %v, want %v", got, tt.want)
			}
		})
	}
}
