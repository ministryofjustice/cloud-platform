#!/usr/bin/env ruby

# Consume the incident log file and output incident performance metrics
#
#   Usage (from inside the `runbooks/source` directory):
#
#     ../bin/mean_time_to_repair.rb
#

require "date"

INCIDENT_LOG = "incident-log.html.md.erb"

QUARTER_REGEX = %r[^## Q\d \d\d\d\d ]
INCIDENT_REGEX = %r[^### Incident on ]
TIME_TO_REPAIR_REGEX = %r[- ..Time to repair..: (\d+h \d+m)]
TIME_TO_RESOLVE_REGEX = %r[- ..Time to resolve..: (\d+h \d+m)]

def main
  puts parse_incident_log
    .map { |quarter, incidents| Quarter.new(quarter, incidents) }
    .sort { |a, b| a.title <=> b.title }
    .map(&:report)
end

# Turn the incident log into a hash of:
#
#   { [quarter label] => [ list of incidents in quarter ], ... }
#
# ...where each 'incident' is a hash: { time_to_repair: "Xh Ym", time_to_resolve: "Xh Ym" }
#
def parse_incident_log
  data = {}
  current_quarter = nil
  current_incident = nil

  IO.foreach(INCIDENT_LOG) do |l|
    line = l.chomp

    case line
    when QUARTER_REGEX
      if current_quarter # i.e. this isn't the first quarter marker in the file
        # We just reached the next quarter, so add the current incident to the current quarter
        data[current_quarter].push(current_incident) unless current_incident.nil?
      end

      # initialise a new 'quarter' hash
      current_quarter = line
      data[current_quarter] = []
    when INCIDENT_REGEX
      # We reached an incident marker. Finish off this incident and start a new hash.
      data[current_quarter].push(current_incident) unless current_incident.nil?
      current_incident = {}
    when TIME_TO_REPAIR_REGEX
      # We've found the time_to_repair line for this incident
      m = TIME_TO_REPAIR_REGEX.match(line)
      current_incident[:time_to_repair] = m[1]
    when TIME_TO_RESOLVE_REGEX
      # We've found the time_to_resolve line for this incident
      m = TIME_TO_RESOLVE_REGEX.match(line)
      current_incident[:time_to_resolve] = m[1]
    end
  end

  # Ensure we handle the last incident in the file
  data[current_quarter].push(current_incident) unless current_incident == {}

  data
end

class Quarter
  attr_reader :title, :incidents

  def initialize(title, incidents)
    @title = title
    @incidents = incidents
      .reject { |i| i == {} }
      .map { |i| Incident.new(i) }
  end

  def report
    <<~EOF
    #{title}
      Incidents: #{incidents.length}
      Mean time to repair: #{mean_time_to_repair}
      Mean time to resolve: #{mean_time_to_resolve}
    EOF
  end

  private

  def mean_time_to_repair
    sum = incidents.map(&:time_to_repair).sum
    hours_and_minutes(sum / incidents.length)
  end

  def mean_time_to_resolve
    sum = incidents.map(&:time_to_resolve).sum
    hours_and_minutes(sum / incidents.length)
  end

  def hours_and_minutes(seconds)
    hours = seconds / 3600
    seconds = seconds % 3600
    minutes = seconds / 60
    "#{hours}h #{minutes}m"
  end
end


class Incident
  def initialize(params)
    @time_to_repair = params[:time_to_repair]
    @time_to_resolve = params[:time_to_resolve]
  end

  def time_to_repair
    to_seconds(@time_to_repair || @time_to_resolve) # In case one is missing
  end

  def time_to_resolve
    to_seconds(@time_to_resolve || @time_to_repair) # In case one is missing
  end

  private

  # "1h 5m" => 3900 (seconds)
  def to_seconds(str)
    hours, minutes = str.split(" ").map(&:to_i)
    (hours * 3600) + (minutes * 60)
  end
end

main

