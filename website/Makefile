all: exec

install:
	gpg --keyserver hkp://keys.gnupg.net --recv-keys 409B6B1796C275462A1703113804BB82D39DC0E3
	\curl -sSL https://get.rvm.io | bash -s stable --ruby
	@sh -c ~/.rvm/scripts/rvm
	gem install bundler

exec: 
	bundle
	bundle exec middleman build
