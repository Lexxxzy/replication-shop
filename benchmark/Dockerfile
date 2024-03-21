FROM mcr.microsoft.com/dotnet/sdk:8.0 AS builder
WORKDIR /app

COPY benchmark/ShopClient.csproj ./
RUN dotnet restore

COPY benchmark ./
RUN dotnet publish -c Release -o out
