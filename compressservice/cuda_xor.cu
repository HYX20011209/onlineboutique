#include <cuda_runtime.h>
#include <stdint.h>
#include <stdlib.h>

__global__ void xor_kernel(const uint8_t* in, uint8_t* out, size_t n)
{
  size_t i = blockIdx.x * blockDim.x + threadIdx.x;
  if (i < n) out[i] = in[i] ^ 0x5A;
}

extern "C"
size_t dummy_gpu_xor(const void* src, size_t len, char** dst_out)
{
  const uint8_t* h_src = static_cast<const uint8_t*>(src);
  uint8_t* d_src;  cudaMalloc(&d_src, len);
  uint8_t* d_dst;  cudaMalloc(&d_dst, len);

  cudaMemcpy(d_src, h_src, len, cudaMemcpyHostToDevice);

  int block = 256;
  int grid  = (len + block - 1) / block;
  xor_kernel<<<grid, block>>>(d_src, d_dst, len);

  uint8_t* h_dst = (uint8_t*)malloc(len);
  cudaMemcpy(h_dst, d_dst, len, cudaMemcpyDeviceToHost);

  cudaFree(d_src);  cudaFree(d_dst);
  *dst_out = reinterpret_cast<char*>(h_dst);
  return len;
}
