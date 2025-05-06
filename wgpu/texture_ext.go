package wgpu

func (p *Texture) AsImageCopy() *TexelCopyTextureInfo {
	return &TexelCopyTextureInfo{
		Texture:  p,
		MipLevel: 0,
		Origin:   Origin3D{},
		Aspect:   TextureAspectAll,
	}
}
