# build application
FROM golang:1.18 as builder

RUN mkdir /work
WORKDIR /work
COPY ./ /work
RUN go build -o /work/app /work/src

# build Docker image
FROM public.ecr.aws/lambda/go:1

COPY --from=builder /work/app ${LAMBDA_TASK_ROOT}
CMD ["app"]