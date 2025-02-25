package config

import (
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Watcher 配置监控器
type Watcher struct {
	viper    *viper.Viper
	onChange func() error
	mu       sync.RWMutex
	done     chan struct{}
}

// NewWatcher 创建配置监控器
func NewWatcher(configFile string, onChange func() error) (*Watcher, error) {
	v := viper.New()
	v.SetConfigFile(configFile)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	w := &Watcher{
		viper:    v,
		onChange: onChange,
		done:     make(chan struct{}),
	}

	// 启动文件监控
	go w.watch()

	return w, nil
}

// watch 监控配置文件变化
func (w *Watcher) watch() {
	w.viper.WatchConfig()
	w.viper.OnConfigChange(func(e fsnotify.Event) {
		// 延迟处理，避免文件系统事件风暴
		time.Sleep(100 * time.Millisecond)

		w.mu.Lock()
		defer w.mu.Unlock()

		if err := w.onChange(); err != nil {
			// TODO: 处理错误，可能需要回滚配置
			return
		}
	})
}

// Stop 停止监控
func (w *Watcher) Stop() {
	close(w.done)
}
