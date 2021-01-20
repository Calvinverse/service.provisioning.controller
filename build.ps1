[CmdletBinding()]
param(
    [string] $dockerTags = '',
    [switch] $Direct = $false
)

. $(Join-Path $PSScriptRoot 'scripts/utilities.ps1')

$revision = Get-Revision
Write-Output "Using revision: '$revision'"

$version = Get-Version -workspaceDirectory $PSScriptRoot
Write-Output "Using version: '$version'"

$date = Get-Date -UFormat '%Y-%m-%dT%T'
Write-Output "Using date: '$date'"

if ($Direct)
{
    Write-Output "Building locally ..."
    New-LocalBuild -date $date -sha1 $revision -version $version -workspaceDirectory $PSScriptRoot
}
else
{
    Write-Output "Building container ..."
    New-Container -date $date -sha1 $revision -version $version -dockerTags $($dockerTags.Split(',', [System.StringSplitOptions]::RemoveEmptyEntries))
}
