MESSAGE = "Update status.json"

system('git config --global user.email "githubactions@example.com"')
system('git config --global user.name "GitHub Actions"')
system("git add .")

latest_commit_message = `git log -1 --pretty=%B`

puts "Last commit: '#{latest_commit_message.strip}'"

if latest_commit_message.include? MESSAGE
  puts "Updating last update commit"
  # also update the date to the latest check time
  system('git commit --amend --no-edit  --date="$(date -R)"')
  system("git push -f")
else
  puts "Creating new update commit..."
  system("git commit -m '#{MESSAGE}'")
  system("git push")
end
