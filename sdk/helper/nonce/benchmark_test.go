package nonce

import (
	"runtime"
    "testing"
	"time"

    "github.com/stretchr/testify/require"
)

const benchValidity = 5*time.Second
const logMemStats = true

func benchWrapper(helper func(*testing.B, NonceService), b *testing.B, s NonceService) {
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)
	helper(b, s)
	runtime.ReadMemStats(&m2)

	if logMemStats {
		b.Logf("in-use memory size:  %v", m2.Alloc - m1.Alloc)
		b.Logf("total alloc size:    %v", m2.TotalAlloc - m1.TotalAlloc)
		b.Logf("in-use memory count: %v", (m2.Mallocs - m2.Frees) - (m1.Mallocs - m1.Frees))
		b.Logf("total alloc count:   %v", m2.Mallocs - m1.Mallocs)
	}
	b.Logf("Tidy output: %v", s.Tidy())
}

func BenchmarkEncryptedNonceServiceGet(b *testing.B) {
	s, err := newEncryptedNonceService(benchValidity)
    require.NoError(b, err)
	benchWrapper(benchGet, b, s)
}

func BenchmarkSyncMapNonceServiceGet(b *testing.B) {
    s, err := newSyncMapNonceService(benchValidity)
    require.NoError(b, err)
	benchWrapper(benchGet, b, s)
}

func benchGet(b *testing.B, s NonceService) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		token, _, err := s.Get()
		require.NoError(b, err)
		require.NotEmpty(b, token)
	}
}

func BenchmarkEncryptedNonceServiceGetRedeem(b *testing.B) {
	s, err := newEncryptedNonceService(benchValidity)
    require.NoError(b, err)
	benchWrapper(benchGetRedeem, b, s)
}

func BenchmarkSyncMapNonceServiceGetRedeem(b *testing.B) {
    s, err := newSyncMapNonceService(benchValidity)
    require.NoError(b, err)
	benchWrapper(benchGetRedeem, b, s)
}

func benchGetRedeem(b *testing.B, s NonceService) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		token, _, err := s.Get()
		require.NoError(b, err)
		require.NotEmpty(b, token)
		ok := s.Redeem(token)
		require.True(b, ok)
		ok = s.Redeem(token)
		require.False(b, ok)
	}
}

func BenchmarkEncryptedNonceServiceGetRedeemTidy(b *testing.B) {
	s, err := newEncryptedNonceService(benchValidity)
    require.NoError(b, err)
	benchWrapper(benchGetRedeemTidy, b, s)
}

func BenchmarkSyncMapNonceServiceGetRedeemTidy(b *testing.B) {
    s, err := newSyncMapNonceService(benchValidity)
    require.NoError(b, err)
	benchWrapper(benchGetRedeemTidy, b, s)
}

func benchGetRedeemTidy(b *testing.B, s NonceService) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		token, _, err := s.Get()
		require.NoError(b, err)
		require.NotEmpty(b, token)
		ok := s.Redeem(token)
		require.True(b, ok)
		s.Tidy()
	}
}

func BenchmarkEncryptedNonceServiceSequentialTidy(b *testing.B) {
	s, err := newEncryptedNonceService(benchValidity)
    require.NoError(b, err)
	benchWrapper(benchGetRedeemSequentialTidy, b, s)
}

func BenchmarkSyncMapNonceServiceSequentialTidy(b *testing.B) {
    s, err := newSyncMapNonceService(benchValidity)
    require.NoError(b, err)
	benchWrapper(benchGetRedeemSequentialTidy, b, s)
}

func benchGetRedeemSequentialTidy(b *testing.B, s NonceService) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		token, _, err := s.Get()
		require.NoError(b, err)
		require.NotEmpty(b, token)
		ok := s.Redeem(token)
		require.True(b, ok)
	}

	s.Tidy()
}

func BenchmarkEncryptedNonceServiceRandomTidy(b *testing.B) {
	s, err := newEncryptedNonceService(benchValidity)
    require.NoError(b, err)
	benchWrapper(benchGetRedeemRandomTidy, b, s)
}

func BenchmarkSyncMapNonceServiceRandomTidy(b *testing.B) {
    s, err := newSyncMapNonceService(benchValidity)
    require.NoError(b, err)
	benchWrapper(benchGetRedeemRandomTidy, b, s)
}

func benchGetRedeemRandomTidy(b *testing.B, s NonceService) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		token, _, err := s.Get()
		require.NoError(b, err)
		require.NotEmpty(b, token)
		if (i % 3) == 1 {
			ok := s.Redeem(token)
			require.True(b, ok)
		}
	}

	s.Tidy()
}

func BenchmarkEncryptedNonceServiceParallelGet(b *testing.B) {
    s, err := newEncryptedNonceService(benchValidity)
    require.NoError(b, err)
    benchWrapper(benchGetParallelGet, b, s)
}

func BenchmarkSyncMapNonceServiceParallelGet(b *testing.B) {
    s, err := newSyncMapNonceService(benchValidity)
    require.NoError(b, err)
    benchWrapper(benchGetParallelGet, b, s)
}

func benchGetParallelGet(b *testing.B, s NonceService) {
    b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
	        token, _, err := s.Get()
    	    require.NoError(b, err)
        	require.NotEmpty(b, token)
		}
    })
}

func BenchmarkEncryptedNonceServiceParallelGetRedeem(b *testing.B) {
    s, err := newEncryptedNonceService(benchValidity)
    require.NoError(b, err)
    benchWrapper(benchGetRedeemParallelGetRedeem, b, s)
}

func BenchmarkSyncMapNonceServiceParallelGetRedeem(b *testing.B) {
    s, err := newSyncMapNonceService(benchValidity)
    require.NoError(b, err)
    benchWrapper(benchGetRedeemParallelGetRedeem, b, s)
}

func benchGetRedeemParallelGetRedeem(b *testing.B, s NonceService) {
    b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
	        token, _, err := s.Get()
    	    require.NoError(b, err)
        	require.NotEmpty(b, token)
			ok := s.Redeem(token)
    	    require.True(b, ok)
		}
    })
}

func BenchmarkEncryptedNonceServiceParallelGetRedeemTidy(b *testing.B) {
    s, err := newEncryptedNonceService(benchValidity)
    require.NoError(b, err)
    benchWrapper(benchGetRedeemParallelGetRedeemTidy, b, s)
}

func BenchmarkSyncMapNonceServiceParallelGetRedeemTidy(b *testing.B) {
    s, err := newSyncMapNonceService(benchValidity)
    require.NoError(b, err)
    benchWrapper(benchGetRedeemParallelGetRedeemTidy, b, s)
}

func benchGetRedeemParallelGetRedeemTidy(b *testing.B, s NonceService) {
    b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
	        token, _, err := s.Get()
    	    require.NoError(b, err)
        	require.NotEmpty(b, token)
			ok := s.Redeem(token)
    	    require.True(b, ok)
			s.Tidy()
		}
    })
}

func BenchmarkEncryptedNonceServiceParallelTidy(b *testing.B) {
    s, err := newEncryptedNonceService(benchValidity)
    require.NoError(b, err)
    benchWrapper(benchParallelTidy, b, s)
}

func BenchmarkSyncMapNonceServiceParallelTidy(b *testing.B) {
    s, err := newSyncMapNonceService(benchValidity)
    require.NoError(b, err)
    benchWrapper(benchParallelTidy, b, s)
}

func benchParallelTidy(b *testing.B, s NonceService) {
    b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
	        token, _, err := s.Get()
    	    require.NoError(b, err)
        	require.NotEmpty(b, token)
			ok := s.Redeem(token)
    	    require.True(b, ok)
		}
    })

	b.StopTimer()
	time.Sleep(2*time.Second + benchValidity)
	runtime.GC()
	b.StartTimer()
	s.Tidy()
}
