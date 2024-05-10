package drm_lease

//go:generate go run github.com/rajveermalviya/go-wayland/cmd/go-wayland-scanner -pkg drm_lease -prefix wp -suffix v1 -o drm_lease.go -i https://gitlab.freedesktop.org/wayland/wayland-protocols/-/raw/1.36/staging/drm-lease/drm-lease-v1.xml
