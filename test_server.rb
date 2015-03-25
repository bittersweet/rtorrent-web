require "xmlrpc/server"

class Service
  def time
    "markie"
  end

  def upcase(s)
    s.upcase
  end

  def sum(x, y)
    x + y
  end

  def print(first, second)
    puts first, second
    "#{first}, #{second}"
  end

  def array(first, second)
    puts first, second
    [[123123213]]
  end

  def version
    {"Version" => "1",
     "version" => "2"}
  end

  def error
    raise XMLRPC::FaultException.new(500, "Server error")
  end
end

puts 'here'
server = XMLRPC::Server.new 5001, 'localhost'
server.add_handler "service", Service.new
server.serve
