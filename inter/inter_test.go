package inter_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cloakwiss/ntdocs/inter"
)

func TestDecomp(t *testing.T) {
	var data string = `H4QHABwHbkzNbiS+YYQks9vr9tOaY9L0HaxiMsGXUPE0FQuF/GIq8v9/bAIVEiEmuM8Q89Opsu3L
VLKUsBACAw9Z+acoLOa1nX9gVClDKF46H8q4Ico4BYo3hMqBRVAWX245VOT4Jkoiy67r6e4uOSml
ZntJXHn9R1g7UKUeum6w8GQMqb1vbllm0jaPJOgpTtMWqLABaOV6JrEB4orZ0hbqvuX0S5AO05oz
qhrjx5X6p94bSmKOh2JSAcHeLFAXEatVZDIT6XRpsjGzbpdQA4zt+nET4O3yD+ygmqEHER5tyQ78
51eUZz+QvbcOtFU89miQANpG46Gi0zi50i9t/fpeOeJLaUyZi/4KPPplZwbL1Oov76+kq7CyDIPv
mBVgmRom5qcsNz8S0amoR1C/Btcr9j2xxTxVPkdVuaKlInG3LFZ9p9jE97gxTnYCsFtMsjh4OErq
1FEdaQl71I2TOvnxVXNAJlhnQJDxf0ynTkm0uoWKta7eoJpueUcMVT2FSDf2I5ljXj3WnAktYw/A
5fsnaTBL8sHO4k1igF9IuyyTdtFV6xIZCBVQm6OkD5tGCc5KKSdTbDgtZdTslxfrmNzUK+JVltww
a6u7At5Vu1iJ1Yn2ULk+X570p3au7EO6jgGiFRLblHRuF8t9e3fWF7tEbW/SltsjdMfPTV79RGby
jvEafnWqnHzg/DxHNnMec8IT+N/jQWMcxvIp/9KJQoWXb9DrRvSC64NpvdCroUI7M34di1IiGflk
jBAODlZvDrQ65a9rzl0C8Js7jy5sVgP++xeJFrhfFN+3PD6fMRMXwxRLl/2iUZJzsPofDHrYw75m
BAFaHcPBAn1fXz0=`

	combined := strings.Join(strings.Split(data, "\n"), "")
	re, err := inter.GetDecompressed(combined)
	if err == nil {
		fmt.Println(string(re))
	}
}
