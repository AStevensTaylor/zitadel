---
title: Connect with Auth0
---

This guide shows how to enable login with ZITADEL on Auth0.

It covers how to:

- create and configure the application in your project
- create and configure the connection in your Auth0 tenant

Prerequisites:

- existing ZITADEL organization, if not present follow [this guide](../../guides/basics/get-started#trying-out-zitadel-on-zitadelch)
- existing project, if not present follow the first 3 steps [here](../../guides/basics/projects#exercise---create-a-simple-project)
- existing Auth0 tenant as described [here](https://auth0.com/docs/get-started/auth0-overview/create-tenants)

> We have to switch between ZITADEL and a Auth0. If the headings begin with "ZITADEL" switch to the ZITADEL Console and if the headings start with Auth0 please switch to the Auth0 GUI.

## **Auth0**: Create a new connection

In Authentication > Enterprise

1. Press the "+" button right to "OpenID Connect"  
  ![Create new connection](/img/oidc/auth0/auth0-create-app.png)
2. Set a connection name for example "ZITADEL"
3. The issuer url is `https://<YOUR_DOMAIN>/.well-known/openid-configuration`
4. Copy the callback URL (ending with `/login/callback`)

The configuration should look like this:

![initial connection configuration](/img/oidc/auth0/auth0-init-app.png)

Next we have to switch to the ZITADEL console.

## **ZITADEL**: Create the application

First of all we create the application in your project.

import CreateApp from "./application/application.mdx";

<CreateApp appType="web" authType="code" appName="Auth0" redirectURI="https://<TENANT>.<REGION>.auth0.com/login/callback"/>

## **Auth0**: Connect ZITADEL

1. Copy the client id from ZITADEL and past it into the **Client ID** field
2. Copy the client secret from ZITADEL and past it into the **Client Secret** field
   ![full configuration](/img/oidc/auth0/auth0-full.png)
3. click Create
4. To verify the connection go to the "Applications" tab and enable the Default App
   [enable app](/img/oidc/auth0/auth0-enable-app.png)
5. Click "Back to OpenID Connect"
6. Click on the "..." button right to the newly created connection and click "Try"
   ![click try](/img/oidc/auth0/auth0-try.png)
7. ZITADEL should open on a new tab and you can enter your login information
8. After you logged in you should see the following:
   ![full configuration](/img/oidc/auth0/auth0-works.png)
