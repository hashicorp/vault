set :base_url, "https://www.vaultproject.io/"

activate :hashicorp do |h|
  h.name         = "vault"
  h.version      = "0.6.5"
  h.github_slug  = "hashicorp/vault"
  h.website_root = "website"
end

helpers do
  # Get the title for the page.
  #
  # @param [Middleman::Page] page
  #
  # @return [String]
  def title_for(page)
    if page && page.data.page_title
      return "#{page.data.page_title} - Vault by HashiCorp"
    end

     "Vault by HashiCorp"
   end

  # Get the description for the page
  #
  # @param [Middleman::Page] page
  #
  # @return [String]
  def description_for(page)
    return escape_html(page.data.description || "")
  end

  # This helps by setting the "active" class for sidebar nav elements
  # if the YAML frontmatter matches the expected value.
  def sidebar_current(expected)
    current = current_page.data.sidebar_current || ""
    if current.start_with?(expected)
      return " class=\"active\""
    else
      return ""
    end
  end

  # Returns the id for this page.
  # @return [String]
  def body_id_for(page)
    if name = page.data.sidebar_current && !name.blank?
      return "page-#{name.strip}"
    end
    return "page-home"
  end

  # Returns the list of classes for this page.
  # @return [String]
  def body_classes_for(page)
    classes = []

    if page && page.data.layout
      classes << "layout-#{page.data.layout}"
    end

    return classes.join(" ")
  end
end
