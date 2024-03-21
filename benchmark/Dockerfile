FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build-env
WORKDIR /app

COPY benchmark/ShopClient.csproj ./
RUN dotnet restore

COPY benchmark ./
RUN dotnet publish -c Release -o out

FROM mcr.microsoft.com/dotnet/sdk:8.0 AS benchmark
WORKDIR /app
COPY --from=build-env /app .
COPY --from=build-env /app/out .
RUN ls