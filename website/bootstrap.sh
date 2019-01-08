echo "==> installing bundler and middleman"
gem install bundler middleman --no-ri --no-rdoc
echo "==> installing ruby dependencies"
bundle
echo "==> installing node dependencies"
cd assets && npm install
