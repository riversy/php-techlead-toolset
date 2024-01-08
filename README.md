# php-techlead-toolset
PHP - Technical Lead toolset


## Usage

### Vendor Compare

This command allow to compare two PHP vendor folders, one is the original and othe one is a clean system, like Magento, Pimcore, Akeneo etc. 
Outcome is the list of custom packages installed over the natice system.

```bash
vendcmp -p /Users/riversy/Workspace/sample-2.4.3-p3 -s /Users/riversy/Workspace/sample-2.4.3-p3/magento
venddiff -d /Users/riversy/Workspace/sample-2.4.3-p3/version_4.0.0.diff
```


