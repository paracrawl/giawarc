
echo "package nbp"
echo

echo "var NBP = map[string]map[string]bool{"
for f in nonbreaking_prefix.*; do
    lang=`echo $f | sed 's/^nonbreaking_prefix\.//'`
    echo "    \"${lang}\": {"
    < $f sed -e '/#NUMERIC_ONLY#/d' -e 's/#.*//' -e 's/ *$//' -e '/^$/d' -e 's/^/        "/' -e 's/$/": true,/'
    echo "    },"
done
echo "}"

echo

echo "var NUM = map[string]map[string]bool{"
for f in nonbreaking_prefix.*; do
    lang=`echo $f | sed 's/^nonbreaking_prefix\.//'`
    echo "    \"${lang}\": {"
    < $f grep '#NUMERIC_ONLY#' | sed -e 's/#.*//' -e 's/ *$//' -e '/^$/d' -e 's/^/        "/' -e 's/$/": true,/'
    echo "    },"
done
echo "}"
