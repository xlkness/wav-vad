#ifndef _WIN32
	#define WEBRTC_POSIX
#endif

#include "cgo/common_audio/signal_processing/complex_bit_reverse.c"
#include "cgo/common_audio/signal_processing/complex_fft.c"
#include "cgo/common_audio/signal_processing/cross_correlation.c"
#include "cgo/common_audio/signal_processing/division_operations.c"
#include "cgo/common_audio/signal_processing/downsample_fast.c"
#include "cgo/common_audio/signal_processing/energy.c"
#include "cgo/common_audio/signal_processing/get_scaling_square.c"
#include "cgo/common_audio/signal_processing/min_max_operations.c"
#include "cgo/common_audio/signal_processing/real_fft.c"
#include "cgo/common_audio/signal_processing/resample_48khz.c"
#include "cgo/common_audio/signal_processing/resample_by_2_internal.c"
#include "cgo/common_audio/signal_processing/resample_fractional.c"
#include "cgo/common_audio/signal_processing/spl_init.c"
#include "cgo/common_audio/signal_processing/vector_scaling_operations.c"
#include "cgo/common_audio/vad/vad_core.c"
#include "cgo/common_audio/vad/vad_filterbank.c"
#include "cgo/common_audio/vad/vad_gmm.c"
#include "cgo/common_audio/vad/vad_sp.c"
#include "cgo/common_audio/vad/webrtc_vad.c"
