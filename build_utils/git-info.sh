commit_seql=$(git rev-parse HEAD)
commit_branch=$(git rev-parse --abbrev-ref HEAD)
commit_detail=$(git log -1)
commit_detail=${commit_detail//"\\"/"\\\\"}

echo "branch:"$commit_seql
echo "commit:"$commit_branch
echo "detail:"$commit_detail

str="
package main;
     func GitInfo() string { return \`_commit_seql:\\t $commit_seql\\n_commit_branch:\\t $commit_branch\\n_commit_detail:\\t $commit_detail\\n  \`;}
"



echo $str > ./git.go
