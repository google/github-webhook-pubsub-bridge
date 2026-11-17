> [!WARNING]
> **Deprecated: this repository is no longer maintained and will be archived
> (made read-only).** Don't deploy it or build on it.
>
> Google now ingests GitHub webhook events through the GitHub Events Platform
> (GEP). Googlers should use GEP (go/github-events-platform) instead of this
> bridge. The open-source code behind GEP is
> [abcxyz/github-metrics-aggregator](https://github.com/abcxyz/github-metrics-aggregator).
>
> Earlier versions of this code hardcoded a webhook secret. Don't reuse it.
> If you deployed this code, rotate your GitHub webhook secret. The final
> version reads the secret from the `GITHUB_WEBHOOK_SECRET` environment
> variable and rejects every request if it isn't set.

# GitHub Webhook to Google Cloud PubSub Bridge

This does exactly what it says it does.

## Source Code Headers

Every file containing source code must include copyright and license
information. This includes any JS/CSS files that you might be serving out to
browsers. (This is to help well-intentioned people avoid accidental copying that
doesn't comply with the license.)

Apache header:

    Copyright 2022 Google LLC

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

        https://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
