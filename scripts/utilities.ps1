
$ErrorActionPreference = 'Stop'
function Get-Revision
{
    [CmdletBinding()]
    param()

    return (git rev-parse HEAD).Trim()
}

function Get-Version
{
    [CmdletBinding()]
    param(
        [string] $workspaceDirectory = $PSScriptRoot
    )

    $tempDir = Join-Path $workspaceDirectory 'temp'
    $gitVersion = Join-Path (Join-Path (Join-Path $tempDir 'GitVersion.CommandLine') 'tools') 'gitversion.exe'
    if (-not (Test-Path $gitVersion))
    {
        if (-not (Test-Path $tempDir))
        {
            New-Item -Path $tempDir -ItemType Directory | Out-Null
        }

        $nuget = Join-Path $tempDir 'nuget.exe'
        if (-not (Test-Path $nuget))
        {
            Invoke-WebRequest -Uri 'https://dist.nuget.org/win-x86-commandline/latest/nuget.exe' -OutFile $nuget
        }

        & nuget install GitVersion.CommandLine -ExcludeVersion -OutputDirectory $tempDir -NonInteractive -Source https://api.nuget.org/v3/index.json | Out-Null
    }

    return & $gitVersion /output json /showvariable SemVer
}

function New-Container
{
    [CmdletBinding()]
    param (
        [string] $sha1,
        [string] $version,
        [string] $date,
        [string[]] $dockerTags = @()
    )

    Write-Output "Building Docker container with build arguments:"
    Write-Output "NOW = $date"
    Write-Output "REVISION = $sha1"
    Write-Output "VERSION = $version"

    $command = "docker build"
    $command += " --force-rm"
    $command += " --build-arg NOW=$date"
    $command += " --build-arg REVISION=$sha1"
    $command += " --build-arg VERSION=$version"
    $command += " --file ./build/package/server/dockerfile"

    if ($dockerTags.Length -gt 0)
    {
        Write-Output "$($dockerTags.Length)"
        foreach($tag in $dockerTags)
        {
            Write-Output "Tagging with: $($tag)/service-provisioning-controller:$version"
            $command += " --tag $($tag)/service-provisioning-controller:$version"
        }
    }
    else
    {
        $command += " --tag service-provisioning-controller:$version"
    }

    $command += " ."

    Write-Output "Invoking: $command"
    Invoke-Expression -Command $command
}

function New-LocalBuild
{
    [CmdletBinding()]
    param (
        [string] $sha1,
        [string] $version,
        [string] $date,
        [string] $workspaceDirectory = $PSScriptRoot
    )

    $outputDir = './bin'
    $absoluteOutputDir = [System.IO.Path]::GetFullPath($(Join-Path $workspaceDirectory $outputDir))
    if (-not (Test-Path $absoluteOutputDir))
    {
        New-Item -Path $absoluteOutputDir -ItemType Directory | Out-Null
    }

    Copy-Item -Path (Join-Path $workspaceDirectory "configs" "*") -Destination $absoluteOutputDir -Force

    & swag init --parseInternal --output ./api --generalInfo ./internal/cmd/serve.go

    $docDirectory = Join-Path $absoluteOutputDir 'api'
    if (-not (Test-Path $docDirectory))
    {
        New-Item -Path $docDirectory -ItemType Directory | Out-Null
    }

    Copy-Item -Path (Join-Path $workspaceDirectory 'api' '*') -Destination $docDirectory -Force

    $configPath = Join-Path $absoluteOutputDir 'config.yaml'
    Add-Content -Path $configPath -Value 'doc:'
    Add-Content -Path $configPath -Value "  path: $docDirectory"

    Add-Content -Path $configPath -Value 'service:'
    Add-Content -Path $configPath -Value '  port: 8080'

    go build -a -installsuffix cgo -v -ldflags="-X github.com/calvinverse/service.provisioning.controller/internal/info.sha1=$sha1 -X github.com/calvinverse/service.provisioning.controller/internal/info.buildTime=$date -X github.com/calvinverse/service.provisioning.controller/internal/info.version=$version" -o $outputDir/controller.exe ./cmd

    go test -cover -coverprofile="$outputDir/coverage.log" -v ./... ./cmd
    go tool cover -html="$outputDir/coverage.log" -o "$outputDir/coverage.html"
}
