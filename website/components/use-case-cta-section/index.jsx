export default function UseCaseCtaSection() {
  return (
    <section className="g-section-block g-cta-section">
      <div>
        <h2>Ready to get started?</h2>
        <a className="g-btn white download" href="/downloads.html">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="20"
            height="22"
            viewBox="0 0 20 22"
          >
            <path d="M9.292 15.706a1 1 0 0 0 1.416 0l3.999-3.999a1 1 0 1 0-1.414-1.414L11 12.586V1a1 1 0 1 0-2 0v11.586l-2.293-2.293a1 1 0 1 0-1.414 1.414l3.999 3.999zM20 16v3c0 1.654-1.346 3-3 3H3c-1.654 0-3-1.346-3-3v-3a1 1 0 1 1 2 0v3c0 .551.448 1 1 1h14c.552 0 1-.449 1-1v-3a1 1 0 1 1 2 0z"></path>
          </svg>
          Download
        </a>
        <a className="g-btn white-outline" href="/docs">
          Explore Docs
        </a>
      </div>
    </section>
  )
}
