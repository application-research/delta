FROM golang:1.18 AS builder
RUN apt-get update && \
    apt-get install -y wget jq hwloc ocl-icd-opencl-dev git libhwloc-dev pkg-config make && \
    apt-get install -y cargo
WORKDIR /app/
ADD . /app
# RUN git init && git add *
RUN curl https://sh.rustup.rs -sSf | bash -s -- -y
ENV PATH="/root/.cargo/bin:${PATH}"
RUN RUSTFLAGS="-C target-cpu=native -g" FFI_BUILD_FROM_SOURCE=1 FFI_USE_BLST_PORTABLE=1 make

FROM golang:1.18
RUN apt-get update && \
    apt-get install -y hwloc libhwloc-dev ocl-icd-opencl-dev
WORKDIR /root/
COPY --from=builder /app/delta ./
CMD ./delta daemon --repo /root/config/.whypfs --wallet-dir /root/config/wallet
