# NOTE: Changes to this file require rebuilding the docker image.
#       Otherwise, your change will not affect the live site.

require 'govuk_tech_docs'

GovukTechDocs.configure(self)

set :layout, 'custom'
