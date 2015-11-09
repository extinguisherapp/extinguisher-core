require_relative './bundle/bundler/setup'
require 'rack'
require 'rack/singleshot'
require 'serve_ext'
require "socket"


class Rack::Handler::SingleShot
  # process request
  def process
    request = read_request

    status, headers, body = @app.call(request)

    write_response(status, headers, body)
  end
end

class SocketServer

  attr_accessor :path

  def initialize(path:, socket_path:)

    @path = path

    File.unlink(socket_path) if File.exists?(socket_path)
    @socket_path = socket_path

    @app = Rack::Builder.app do

      use Rack::CommonLogger
      use Rack::ShowStatus
      use Rack::ShowExceptions
      run Rack::Cascade.new([
        Serve::RackAdapter.new(path),
        Rack::Directory.new(path)
      ])
    end

    @server ||= UNIXServer.new(@socket_path)
  end

  def start
    while true
      socket = @server.accept
      Rack::Handler::SingleShot.new(@app, socket, socket).process
      socket.close
    end
  end
end

Dir.chdir(ARGV[0])
SocketServer.new(path: ARGV[0], socket_path: File.join(ARGV[0],"html.socket")).start
