git checkout -b nextjs-fe : create branch nextjs-fe
git checkout -b nextjs-fe-reading

git branch -v: see branch local

git branch -vv
: to see all local and remote branches with detail

git push --set-upstream origin nextjs-fe-reading
git push -u origin head: This command will create the branch on the remote if it doesn’t exist yet


git checkout nextjs-fe-reading: jump to nextjs-fe


git checkout nextjs-fe-reading


---------------------------------- merge nrach --------
git checkout nextjs-fe-reading
git merge nextjs-fe-readiding-add-more-clean-architextture
git commit -m "Merged nextjs-fe-readiding-add-more-clean-architextture into nextjs-fe-reading"
git push origin nextjs-fe-reading
-------------------------------- done merge ----------------