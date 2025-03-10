/** @type {import('@docusaurus/types').DocusaurusConfig} */
module.exports = {
  title: "ZITADEL Docs",
  trailingSlash: false,
  url: "https://docs.zitadel.com",
  baseUrl: "/",
  onBrokenLinks: "warn",
  onBrokenMarkdownLinks: "warn",
  favicon: "img/favicon.ico",
  organizationName: "zitadel",
  projectName: "zitadel",
  scripts: [
    {
      src: "/proxy/js/script.js",
      async: true,
      defer: true,
      "data-domain": "docs.zitadel.com",
      "data-api": "/proxy/api/event",
    },
  ],
  themeConfig: {
    zoomSelector: ".markdown :not(em) > img",
    announcementBar: {
      id: 'documentation',
      content:
        'This page contains the documentation for ZITADEL version 2, if you are looking for version 1 please visit <a target="_blank" rel="noopener noreferrer" href="https://docs-v1.zitadel.com">https://docs-v1.zitadel.com</a>',
      backgroundColor: '#fafbfc',
      textColor: '#091E42',
      isCloseable: false,
    },
    navbar: {
      // title: 'ZITADEL',
      logo: {
        alt: "ZITADEL logo",
        src: "img/zitadel-logo-dark.svg",
        srcDark: "img/zitadel-logo-light.svg",
      },
      items: [
        {
          type: "doc",
          label: "Guides",
          docId: "guides/overview",
          position: "left",
        },
        {
          type: "doc",
          label: "Quickstarts",
          docId: "quickstarts/introduction",
          position: "left",
        },
        {
          type: "doc",
          label: "APIs",
          docId: "apis/introduction",
          position: "left",
        },
        {
          type: "doc",
          docId: "concepts/introduction",
          label: "Concepts",
          position: "left",
        },
        {
          type: "doc",
          docId: "manuals/introduction",
          label: "Help",
          position: "left",
        },
        {
          type: "doc",
          docId: "legal/introduction",
          label: "Legal",
          position: "left",
        },
        {
          href: "https://github.com/zitadel/zitadel",
          label: "GitHub",
          position: "right",
        },
      ],
    },
    footer: {
      links: [
        {
          title: "Community",
          items: [
            {
              label: "Chat",
              href: "https://zitadel.com/chat",
            },
            {
              label: "GitHub Discussions",
              href: "https://github.com/zitadel/zitadel/discussions",
            },
            {
              label: "Twitter",
              href: "https://twitter.com/zitadel",
            },
            {
              label: "Linkedin",
              href: "https://www.linkedin.com/company/zitadel/",
            },
            {
              label: "Blog",
              href: "https://zitadel.com/blog",
            },
          ],
        },
        {
          title: "Legal",
          items: [

            {
              label: "Terms and Conditions",
              href: "/docs/legal/terms-of-service",
            },
            {
              label: "Privacy Policy",
              href: "/docs/legal/privacy-policy",
            },
          ],
        },
        {
          title: "About",
          items: [
            {
              label: "Website",
              href: "https://zitadel.com",
            },
            {
              label: "Contact",
              href: "https://zitadel.com/contact/",
            },
            {
              label: "GitHub",
              href: "https://github.com/zitadel",
            },
            {
              label: "Status",
              href: "https://status.zitadel.ch/",
            }
          ],
        },
        
      ],
      copyright: `Copyright © ${new Date().getFullYear()} ZITADEL Docs - Built with Docusaurus.`,
    },
    algolia: {
      appId: "8H6ZKXENLO",
      apiKey: "124fe1c102a184bc6fc70c75dc84f96f",
      indexName: 'zitadel',
      selector: 'div#'
  },
    prism: {
      additionalLanguages: ["csharp", "dart", "groovy"],
    },
  },
  presets: [
    [
      "@docusaurus/preset-classic",
      {
        docs: {
          sidebarPath: require.resolve("./sidebars.js"),
          editUrl: "https://github.com/zitadel/zitadel/edit/v2-alpha/docs/",
          remarkPlugins: [require("mdx-mermaid")],
        },
        theme: {
          customCss: require.resolve("./src/css/custom.css"),
        },
      },
    ],
  ],
  plugins: [require.resolve("plugin-image-zoom")],
};
