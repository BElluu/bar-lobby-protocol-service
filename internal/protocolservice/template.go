package protocolservice

import "html/template"

var pageTemplate = template.Must(template.New("protocol-page").Parse(protocolPageHTML))

const protocolPageHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>BAR Lobby</title>
  <style nonce="{{ .Nonce }}">
    :root {
      color-scheme: dark;
      --bg: #080b0f;
      --panel: rgba(12, 17, 23, 0.82);
      --panel-border: rgba(255, 255, 255, 0.14);
      --text: #f4f7fb;
      --muted: #b5c0cc;
      --accent: #f0c75e;
      --accent-strong: #ffe28a;
      --ink: #111318;
    }

    * {
      box-sizing: border-box;
    }

    body {
      min-height: 100vh;
      margin: 0;
      color: var(--text);
      font-family: Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
      background:
        linear-gradient(90deg, rgba(8, 11, 15, 0.95) 0%, rgba(8, 11, 15, 0.76) 46%, rgba(8, 11, 15, 0.38) 100%),
        url("/assets/background.avif") center / cover fixed,
        var(--bg);
    }

    main {
      min-height: 100vh;
      display: grid;
      align-items: center;
      justify-items: center;
      padding: clamp(24px, 6vw, 72px);
    }

    .panel {
      width: min(100%, 560px);
      padding: clamp(24px, 5vw, 42px);
      border: 1px solid var(--panel-border);
      border-radius: 8px;
      background: var(--panel);
      box-shadow: 0 24px 80px rgba(0, 0, 0, 0.48);
      backdrop-filter: blur(16px);
    }

    .brand {
      width: min(300px, 72vw);
      height: auto;
      display: block;
      margin-bottom: 34px;
    }

    h1 {
      margin: 0 0 14px;
      font-size: clamp(2.1rem, 7vw, 4.5rem);
      line-height: 0.94;
      letter-spacing: 0;
      text-transform: uppercase;
    }

    p {
      margin: 0;
      color: var(--muted);
      font-size: 1.05rem;
      line-height: 1.6;
    }

    .actions {
      display: flex;
      flex-wrap: wrap;
      align-items: center;
      gap: 16px;
      margin-top: 30px;
    }

    .button {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      min-height: 48px;
      padding: 0 22px;
      border-radius: 6px;
      color: var(--ink);
      background: linear-gradient(180deg, var(--accent-strong), var(--accent));
      font-weight: 800;
      text-decoration: none;
      text-transform: uppercase;
      box-shadow: 0 10px 30px rgba(240, 199, 94, 0.2);
    }

    @media (max-width: 560px) {
      main {
        align-items: center;
        padding: 18px;
      }

      .panel {
        padding: 22px;
      }

      .actions {
        align-items: stretch;
        flex-direction: column;
      }

      .button {
        width: 100%;
      }
    }
  </style>
</head>
<body>
  <main>
    <section class="panel" aria-labelledby="title">
      <img class="brand" src="/assets/bar-logo.avif" alt="Beyond All Reason">
      <h1 id="title">BAR Lobby</h1>
      <p>This link is ready for the BAR Lobby desktop app. Your browser may ask for permission before switching applications.</p>

      <div class="actions">
        <a class="button" href="{{ .ProtocolHref }}">Click here if nothing happens</a>
      </div>
    </section>
  </main>

  <script nonce="{{ .Nonce }}">
    window.addEventListener("load", function () {
      window.location.href = "{{ .ProtocolURL }}";
    });
  </script>
</body>
</html>
`
