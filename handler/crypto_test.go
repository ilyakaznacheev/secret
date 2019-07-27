package handler

import (
	"fmt"
	"testing"
)

func Test_crypto(t *testing.T) {
	tests := []struct {
		name     string
		genTimes int
		secret   string
		wantErr  bool
	}{
		{
			name:     "simple",
			genTimes: 10,
			secret:   "test",
			wantErr:  false,
		},

		{
			name:     "not so simple",
			genTimes: 100,
			secret:   "test testy testest",
			wantErr:  false,
		},

		{
			name:     "not simple at all",
			genTimes: 1000,
			secret:   "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec tincidunt justo a eros accumsan, id efficitur metus dapibus. Mauris ac felis nulla. Morbi ultricies metus et mi pretium hendrerit. Nullam vitae nisi leo. Sed arcu erat, imperdiet nec ante et, semper rhoncus eros. Proin sagittis orci erat, nec consequat ipsum lacinia id. Aenean convallis nisl quis pharetra malesuada. Integer sit amet dignissim quam, in mollis lorem. Sed at sapien lectus. Quisque in sem et dolor ullamcorper feugiat nec vitae enim. Suspendisse vitae luctus enim. Duis sit amet metus in felis consequat suscipit. Nam sed urna aliquet, dignissim arcu eget, lobortis nisi. Pellentesque varius nulla lectus, interdum blandit augue dapibus non. In vulputate, metus sed luctus imperdiet, leo orci molestie purus, nec posuere neque mauris non purus.",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		for idx := 0; idx < tt.genTimes; idx++ {
			t.Run(fmt.Sprintf("%s[%d]", tt.name, idx), func(t *testing.T) {
				key := generateKey()

				encSecret, err := encryptSecret(key, tt.secret)
				if (err != nil) != tt.wantErr {
					t.Errorf("encryptSecret() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				decSecret, err := decryptSecret(key, encSecret)
				if (err != nil) != tt.wantErr {
					t.Errorf("decryptSecret() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if err == nil && decSecret != tt.secret {
					t.Errorf("wrong decoded secret %v, want %v", decSecret, tt.secret)
				}
			})

		}
	}
}
